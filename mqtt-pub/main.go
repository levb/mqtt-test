package main

import (
	"encoding/json"
	"os"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	mqttapp "github.com/levb/mqtt-test"
	server "github.com/nats-io/nats-server/v2/server"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "pub [--flags...]",
	Short: "MQTT Publish",
	Run: func(cmd *cobra.Command, args []string) {
		err := run()
		if err != nil {
			panic("Error: " + err.Error())
		}
	},
}

var opts = struct {
	*mqttapp.Options
	Retain bool
	Topic  string
	QOS    int
	N      int
	Size   int
}{
	Options: &mqttapp.Opts,
}

func init() {
	mqttapp.InitCommand(cmd)

	cmd.Flags().BoolVar(&opts.Retain, "retain", false, "Mark each message as retained")
	cmd.Flags().StringVar(&opts.Topic, "topic", mqttapp.DefaultTopic, "MQTT topic")
	cmd.Flags().IntVar(&opts.QOS, "qos", mqttapp.DefaultQOS, "MQTT QOS")
	cmd.Flags().IntVar(&opts.N, "n", 1, "Number of messages to publish")
	cmd.Flags().IntVar(&opts.Size, "size", 0, "Approximate size of each message (pub adds a timestamp)")
}

func main() {
	_ = cmd.Execute()
}

func run() error {
	clientOpts := paho.NewClientOptions().
		SetProtocolVersion(4).
		SetClientID(opts.ClientID).
		SetUsername(opts.Username).
		SetPassword(opts.Password).
		SetStore(paho.NewMemoryStore()).
		SetCleanSession(true)
	for _, s := range opts.Servers {
		clientOpts.AddBroker(s)
	}
	cl := paho.NewClient(clientOpts)

	if t := cl.Connect(); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	defer cl.Disconnect(mqttapp.DisconnectCleanupTimeout)

	// ready to publish
	os.Stdout.Write(mqttapp.READY)

	elapsed := time.Duration(0)
	bc := 0
	for n := 0; n < opts.N; n++ {
		if n > 0 {
			time.Sleep(1 * time.Millisecond)
		}

		// payload always starts with JSON containing timestamp, etc. The JSON
		// is always terminated with a '-', which can not be part of the random
		// fill. payload is then filled to the requested size with random data.
		payload := mqttapp.RandomPayload(opts.Size)
		structuredPayload, _ := json.Marshal(mqttapp.PubValue{
			Seq:       n,
			Timestamp: time.Now().UnixNano(),
		})
		structuredPayload = append(structuredPayload, '\n')
		if len(structuredPayload) > len(payload) {
			payload = structuredPayload
		} else {
			copy(payload, structuredPayload)
		}
		start := time.Now()
		if token := cl.Publish(opts.Topic, byte(opts.QOS), opts.Retain, payload); token.Wait() && token.Error() != nil {
			return token.Error()
		}

		elapsed += time.Since(start)
		bc += mqttapp.LenPublish(opts.Topic, byte(opts.QOS), opts.Retain, payload)
	}

	bb, _ := json.Marshal(server.MQTTBenchmarkResult{
		Ops:   opts.N,
		NS:    elapsed,
		Unit:  "pub",
		Bytes: int64(bc),
	})
	os.Stdout.Write(bb)

	return nil
}

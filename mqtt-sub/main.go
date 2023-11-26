package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	mqttapp "github.com/levb/mqtt-test"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "sub [--flags...]",
	Short: "MQTT Subscribe",
	Run: func(cmd *cobra.Command, args []string) {
		err := run()
		if err != nil {
			log.Fatal("Error: " + err.Error())
		}
	},
}

var opts = struct {
	*mqttapp.Options
	Topic            string
	MatchTopicPrefix string
	QOS              int
	N                int
}{
	Options: &mqttapp.Opts,
}

func init() {
	mqttapp.InitCommand(cmd)

	cmd.Flags().StringVar(&opts.Topic, "topic", mqttapp.DefaultTopic, "MQTT topic, can be a wildcard")
	cmd.Flags().StringVar(&opts.MatchTopicPrefix, "match-prefix", mqttapp.DefaultTopic, "Ignore mesages that don't match the prefix")
	cmd.Flags().IntVar(&opts.QOS, "qos", mqttapp.DefaultQOS, "MQTT QOS")
	cmd.Flags().IntVar(&opts.N, "n", 1, "Number of messages to publish")
}

func main() {
	_ = cmd.Execute()
}

func run() error {
	msgChan := make(chan paho.Message)
	errChan := make(chan error)

	ready := sync.WaitGroup{}
	ready.Add(1)
	clientOpts := paho.NewClientOptions().
		SetProtocolVersion(4).
		SetClientID(opts.ClientID).
		SetUsername(opts.Username).
		SetPassword(opts.Password).
		SetStore(paho.NewMemoryStore()).
		SetAutoAckDisabled(true).
		SetCleanSession(true).
		SetOnConnectHandler(func(cl paho.Client) {
			token := cl.Subscribe(opts.Topic, byte(opts.QOS), func(client paho.Client, msg paho.Message) {
				msg.Ack()
				switch {
				case opts.MatchTopicPrefix != "" && !strings.HasPrefix(msg.Topic(), opts.MatchTopicPrefix):
					// ignore

				case msg.Duplicate():
					// ignore
					errChan <- fmt.Errorf("received unexpected duplicate message")

				case msg.Retained():
					errChan <- fmt.Errorf("received unexpected retained message")

				default:
					msgChan <- msg
				}
			})
			if token.Wait() && token.Error() != nil {
				errChan <- token.Error()
			}
			ready.Done()
		}).
		SetDefaultPublishHandler(func(client paho.Client, msg paho.Message) {
			log.Printf("<>/<> 100\n")
			// <>/<> TODO
		})

	for _, s := range opts.Servers {
		clientOpts.AddBroker(s)
	}
	cl := paho.NewClient(clientOpts)

	if t := cl.Connect(); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	defer cl.Disconnect(mqttapp.DisconnectCleanupTimeout)

	ready.Wait()
	os.Stdout.Write(mqttapp.READY)

	elapsed := time.Duration(0)
	bc := 0
	timeout := time.After(mqttapp.IdleTimeout)
	for n := 0; n < opts.N; {
		select {
		case <-timeout:
			log.Fatalf("Error: timeout waiting for messages")

		case err := <-errChan:
			log.Fatalf("Error: %v", err)

		case msg := <-msgChan:
			v := mqttapp.PubValue{}
			body := msg.Payload()
			if i := bytes.IndexByte(body, '\n'); i != -1 {
				body = body[:i]
			}
			err := json.Unmarshal(body, &v)
			if err != nil {
				log.Fatalf("Error parsing message JSON: %v", err)
			}
			elapsed += time.Since(time.Unix(0, v.Timestamp))
			bc += len(msg.Payload())
			n++
		}
	}

	bb, _ := json.Marshal(mqttapp.MQTTBenchmarkResult{
		Ops:   opts.N,
		NS:    elapsed,
		Unit:  "sub",
		Bytes: int64(bc),
	})
	os.Stdout.Write(bb)

	return nil
}

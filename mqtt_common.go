package mqtt

import (
	"log"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
)

const (
	DefaultServer     = "tcp://localhost:1883"
	DefaultClientID   = "test"
	DefaultUsername   = ""
	DefaultPassword   = ""
	DefaultPubPrefix  = "pub"
	DefaultSubPrefix  = "sub"
	DefaultTopic      = "foo"
	DefaultQOS        = 0
	DefaultQOSTimeout = 60 * time.Second

	IdleTimeout              = 10 * time.Second
	DisconnectCleanupTimeout = 500 // milliseconds
)

var READY = []byte("READY\n")

type Options struct {
	ClientID string
	Servers  []string
	Username string
	Password string
}

type PubValue struct {
	Seq       int   `json:"seq"`
	Timestamp int64 `json:"timestamp"`
}

type MQTTBenchmarkResult struct {
	Ops   int           `json:"ops"`
	NS    time.Duration `json:"ns"`
	Unit  string        `json:"unit"`
	Bytes int64         `json:"bytes"`
}


var Opts Options

func InitCommand(cmd *cobra.Command) {
	mqtt.ERROR = log.New(os.Stderr, "[MQTT ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stderr, "[MQTT CRIT] ", 0)
	mqtt.WARN = log.New(os.Stderr, "[MQTT WARN] ", 0)
	cmd.Flags().StringVar(&Opts.ClientID, "id", DefaultClientID, "MQTT client ID")
	cmd.Flags().StringArrayVarP(&Opts.Servers, "server", "s", []string{DefaultServer}, "MQTT servers endpoint as host:port")
	cmd.Flags().StringVarP(&Opts.Username, "username", "u", DefaultUsername, "MQTT client username (empty if auth disabled)")
	cmd.Flags().StringVarP(&Opts.Password, "password", "p", DefaultPassword, "MQTT client password (empty if auth disabled)")
}

var ch = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@$#%^&*()")

func RandomPayload(sz int) []byte {
	b := make([]byte, sz)
	for i := range b {
		b[i] = ch[rand.Intn(len(ch))]
	}
	return b
}

func LenVarInt(value int) int {
	c := 0
	for ; value > 0; value >>= 7 {
		c++
	}
	return c
}

func LenPublish(topic string, qos byte, retained bool, msg []byte) int {
	// Compute len (will have to add packet id if message is sent as QoS>=1)
	pkLen := 2 + len(topic) + len(msg)
	if qos > 0 {
		pkLen += 2
	}
	return 1 + LenVarInt(pkLen) + pkLen
}

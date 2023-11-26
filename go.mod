module github.com/levb/mqtt-test

go 1.21.3

require (
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
)

replace github.com/nats-io/nats-server/v2 => ../../nats-io/nats-server

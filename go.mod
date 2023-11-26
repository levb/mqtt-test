module github.com/levb/mqtt-test

go 1.21.3

require (
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/nats-io/nats-server/v2 v2.10.5
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.17.3 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.5.3 // indirect
	github.com/nats-io/nkeys v0.4.6 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/time v0.4.0 // indirect
)

replace github.com/nats-io/nats-server/v2 => ../../nats-io/nats-server

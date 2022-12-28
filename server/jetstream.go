package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"main/service"
	"main/storage"
	"time"

	"github.com/nats-io/nats.go"
)

func Createsub(subj string) {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
	opts = setupConnOptions(opts)
	// Connect to NATS
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal("failed stream creation", err)
	}
	memory := storage.GetCache()
	_, err = js.Subscribe(subj, func(msg *nats.Msg) {
		log.Printf("message received on %s\n", subj)

		var m service.Message
		err := json.Unmarshal(msg.Data, &m)

		if err != nil {
			fmt.Println(err)
		} else {
			err = memory.Set(m.Order_uid, m)
			if err != nil {
				log.Println(err)
			}
			storage.AddToDB(m)

		}
	})
	
	return

}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

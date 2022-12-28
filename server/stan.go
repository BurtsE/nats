package server

import (
	"encoding/json"
	"fmt"
	"log"
	"main/service"
	"main/storage"

	"github.com/nats-io/stan.go"
)

func CreateStansub(subj string) {
	log.Println("Creating subscriber")
	sc, _ := stan.Connect("test-cluster", "clientid")
	memory := storage.GetCache()
	_, err := sc.Subscribe(subj, func(msg *stan.Msg) {
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

	if err != nil {
		log.Println(err)
	}
}

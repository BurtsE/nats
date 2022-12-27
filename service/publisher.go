package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type Message struct {
	Order_uid    string `json: "order_uid, omitempty"`
	Track_number string `json: "track_number"`
	Entry        string `json: "entry"`
	Delivery     struct {
		Name    string `json: "name"`
		Phone   string `json: "phone"`
		Zip     string `json: "zip"`
		City    string `json: "city"`
		Address string `json: "address"`
		Region  string `json: "region"`
		Email   string `json: "email"`
	}
	Payment struct {
		Transaction   string `json: "transaction"`
		Request_id    string `json: "request_id"`
		Currency      string `json: "currency"`
		Provider      string `json: "provider"`
		Amount        int    `json: "amount"`
		Payment_dt    int64  `json: "payment_dt"`
		Bank          string `json: "bank"`
		Delivery_cost int    `json: "delivery_cost"`
		Goods_total   int    `json: "goods_total"`
		Custom_fee    int    `json: "custom_fee"`
	}
	Items              []Item
	Locale             string `json: "locale"`
	Internal_signature string `json: "internal_signature"`
	Customer_id        string `json: "customer_id"`
	Delivery_service   string `json: "delivery_service"`
	Shardkey           string `json: "shardkey"`
	Sm_id              int    `json: "sm_id"`
	Date_created       string `json: "date_created"`
	Oof_shard          string `json: "oof_shard"`
}

type Item struct {
	Chrt_id      int    `json: "chrt_id"`
	Track_number string `json: "track_number"`
	Price        int    `json: "price"`
	Rid          string `json: "rid"`
	Name         string `json: "name"`
	Sale         int    `json: "sale"`
	Size         string `json: "size"`
	Total_price  int    `json: "total_price"`
	Nm_id        int    `json: "nm_id"`
	Brand        string `json: "brand"`
	Status       int    `json: "status"`
}

func LaunchPublisher(js nats.JetStreamContext, consumer string) {
	go func() {
		log.Println("Publishing launched")
		log.Println("stream added")
		for false {
			file, _ := os.Open("model.json")
			data, _ := ioutil.ReadAll(file)
			file.Close()
			js.PublishAsync(consumer, data)
			file, _ = os.Open("model2.json")
			data, _ = ioutil.ReadAll(file)
			file.Close()
			js.PublishAsync(consumer, data)
			file, _ = os.Open("model3.json")
			data, _ = ioutil.ReadAll(file)
			file.Close()
			js.PublishAsync(consumer, data)

			log.Println("Publishing finished")
			select {
			case <-js.PublishAsyncComplete():
			case <-time.After(5 * time.Second):
				fmt.Println("Did not resolve in time")
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"math/rand"
	"os"
	"projectL0/internal/entity"
	"time"
)

const (
	streamName          = "ORDERS"
	streamSubjects      = "ORDERS.*"
	subjectNameNewOrder = "ORDERS.NEW"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))

	stream, _ := js.StreamInfo(streamName)

	if stream == nil {
		_, err := js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubjects},
		})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	publishOrder(js)
}

func publishOrder(js nats.JetStreamContext) {
	order, err := getOrder()
	if err != nil {
		log.Fatalln(err.Error())
	}

	r := rand.Intn(10)
	time.Sleep(time.Duration(r) * time.Second)

	orderString, err := json.Marshal(order)
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = js.Publish(subjectNameNewOrder, orderString)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Publisher  =>  order_uid:%s\n", order.OrderUID)
	}
}

func getOrder() (entity.Order, error) {
	rawOrder, _ := os.ReadFile("../model.json")
	var order entity.Order
	err := json.Unmarshal(rawOrder, &order)

	return order, err
}

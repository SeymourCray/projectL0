package jetstream

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"projectL0/internal/entity"
	"projectL0/internal/repository"
)

const (
	subjectName = "ORDERS.NEW"
)

type OrderJetstream struct {
	js   nats.JetStreamContext
	repo repository.OrderRepo
}

func NewOrderJetstream(js nats.JetStreamContext, repo repository.OrderRepo) *OrderJetstream {
	return &OrderJetstream{js: js, repo: repo}
}

func (orderJetstream *OrderJetstream) ConsumeOrders() {
	_, err := orderJetstream.js.Subscribe(subjectName, func(m *nats.Msg) {
		err := m.Ack()

		if err != nil {
			log.Println(err.Error())
		}

		var order entity.Order
		err = json.Unmarshal(m.Data, &order)
		if err != nil {
			log.Println(err.Error())
		}

		orderJetstream.repo.InsertOrder(order)
	})

	if err != nil {
		log.Fatalln("Subscribe failed: " + err.Error())
	}
}

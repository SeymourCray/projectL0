package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"log"
	"projectL0/config"
	"projectL0/internal/controller/http"
	"projectL0/internal/controller/jetstream"
	"projectL0/internal/repository"
)

func Run(cfg config.Config) {
	conn, err := sqlx.Connect("pgx", cfg.PostgresURL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	orderRepository := repository.NewOrderRepository(conn)

	viewsEngine := html.New("../web/views", ".html")

	webApp := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	orderController := http.NewOrderController(webApp, orderRepository)
	orderController.RegisterHandlers()

	nc, _ := nats.Connect(nats.DefaultURL)
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))

	orderJetstream := jetstream.NewOrderJetstream(js, orderRepository)
	orderJetstream.ConsumeOrders()

	port := fmt.Sprintf(":%s", cfg.HttpPort)
	if err := webApp.Listen(port); err != nil {
		log.Fatalln(err)
	}
}

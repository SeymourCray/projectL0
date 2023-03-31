package http

import (
	"github.com/gofiber/fiber/v2"
	"projectL0/internal/repository"
)

type OrderController struct {
	app  *fiber.App
	repo repository.OrderRepo
}

func NewOrderController(app *fiber.App, repo repository.OrderRepo) *OrderController {
	return &OrderController{app: app, repo: repo}
}

func (controller *OrderController) GetOrder(ctx *fiber.Ctx) error {
	orderUID := ctx.Query("order_uid")

	order, _ := controller.repo.GetOrderByID(orderUID)

	return ctx.Render("order", order)
}

func (controller *OrderController) RegisterHandlers() {
	controller.app.Get("/orders", controller.GetOrder)
}

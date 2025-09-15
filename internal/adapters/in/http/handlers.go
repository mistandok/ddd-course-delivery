package http

import (
	"delivery/internal/generated/servers"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) GetCouriers(ctx echo.Context) error {
	// TODO: implement
	return ctx.JSON(200, []servers.Courier{})
}

func (h *Handlers) CreateCourier(ctx echo.Context) error {
	// TODO: implement
	return ctx.NoContent(201)
}

func (h *Handlers) CreateOrder(ctx echo.Context) error {
	// TODO: implement
	return ctx.NoContent(201)
}

func (h *Handlers) GetOrders(ctx echo.Context) error {
	// TODO: implement
	return ctx.JSON(200, []servers.Order{})
}
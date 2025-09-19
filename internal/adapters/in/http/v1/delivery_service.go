package v1

import (
	"log"
	"net/http"

	"delivery/internal/core/application/usecases/commands/create_courier"
	"delivery/internal/core/application/usecases/commands/create_order"
	"delivery/internal/core/application/usecases/queries/get_all_couriers"
	"delivery/internal/core/application/usecases/queries/get_all_uncompleted_orders"
	"delivery/internal/generated/servers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DeliveryService struct {
	getAllCouriersHandler          get_all_couriers.GetAllCouriersHandler
	createCourierHandler           create_courier.CreateCourierHandler
	getAllUncompletedOrdersHandler get_all_uncompleted_orders.GetAllUncompletedOrdersHandler
	createOrderHandler             create_order.CreateOrderHandler
}

func NewDeliveryService(
	getAllCouriersHandler get_all_couriers.GetAllCouriersHandler,
	createCourierHandler create_courier.CreateCourierHandler,
	getAllUncompletedOrdersHandler get_all_uncompleted_orders.GetAllUncompletedOrdersHandler,
	createOrderHandler create_order.CreateOrderHandler,
) *DeliveryService {
	return &DeliveryService{
		getAllCouriersHandler:          getAllCouriersHandler,
		createCourierHandler:           createCourierHandler,
		getAllUncompletedOrdersHandler: getAllUncompletedOrdersHandler,
		createOrderHandler:             createOrderHandler,
	}
}

func (d *DeliveryService) CreateCourier(ctx echo.Context) error {
	var newCourier servers.NewCourier
	if err := ctx.Bind(&newCourier); err != nil {
		return ctx.JSON(http.StatusBadRequest, servers.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	command, err := create_courier.NewCreateCourierCommand(newCourier.Name, int64(newCourier.Speed))
	if err != nil {
		return err
	}

	err = d.createCourierHandler.Handle(ctx.Request().Context(), command)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusCreated)
}

func (d *DeliveryService) GetCouriers(ctx echo.Context) error {
	query := get_all_couriers.NewGetAllCouriersQuery()

	response, err := d.getAllCouriersHandler.Handle(ctx.Request().Context(), query)
	if err != nil {
		return err
	}

	couriers := make([]servers.Courier, len(response.Couriers))
	for i, courierDTO := range response.Couriers {
		couriers[i] = servers.Courier{
			Id:   courierDTO.ID,
			Name: courierDTO.Name,
			Location: servers.Location{
				X: int(courierDTO.Location.X),
				Y: int(courierDTO.Location.Y),
			},
		}
	}

	return ctx.JSON(http.StatusOK, couriers)
}

func (d *DeliveryService) CreateOrder(ctx echo.Context) error {
	orderID := uuid.New()
	command, err := create_order.NewCreateOrderCommand(orderID, "default_street", 1)
	if err != nil {
		return err
	}

	err = d.createOrderHandler.Handle(ctx.Request().Context(), command)
	if err != nil {
		log.Printf("error creating order: %v", err)
		return err
	}

	return ctx.NoContent(http.StatusCreated)
}

func (d *DeliveryService) GetOrders(ctx echo.Context) error {
	query := get_all_uncompleted_orders.NewGetAllUncompletedOrdersQuery()

	response, err := d.getAllUncompletedOrdersHandler.Handle(ctx.Request().Context(), query)
	if err != nil {
		return err
	}

	orders := make([]servers.Order, len(response.Orders))
	for i, orderDTO := range response.Orders {
		orders[i] = servers.Order{
			Id: orderDTO.ID,
			Location: servers.Location{
				X: int(orderDTO.Location.X),
				Y: int(orderDTO.Location.Y),
			},
		}
	}

	return ctx.JSON(http.StatusOK, orders)
}

package services

import (
	aggCourier "delivery/internal/core/domain/model/courier"
	aggOrder "delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"errors"
	"math"
)

type Dispatcher interface {
	Dispatch(order *aggOrder.Order, couriers []*aggCourier.Courier) (*aggCourier.Courier, error)
}

var _ Dispatcher = (*CourierDispatcher)(nil)

type CourierDispatcher struct{}

func NewCourierDispatcher() *CourierDispatcher {
	return &CourierDispatcher{}
}

func (c *CourierDispatcher) Dispatch(order *aggOrder.Order, couriers []*aggCourier.Courier) (*aggCourier.Courier, error) {
	if order == nil {
		return nil, errs.NewValueIsInvalidErrorWithCause("order", errors.New("impossible to dispatch order without order"))
	}

	if !aggOrder.StatusCreated.Equals(order.Status()) {
		return nil, errs.NewValueIsInvalidErrorWithCause("order", errors.New("impossible to dispatch order in status other than created"))
	}

	if len(couriers) == 0 {
		return nil, errs.NewValueIsInvalidErrorWithCause("couriers", errors.New("impossible to dispatch order without couriers"))
	}

	bestCourier, err := c.selectBestCourier(order, couriers)
	if err != nil {
		return nil, err
	}

	err = bestCourier.TakeOrder(order)
	if err != nil {
		return nil, err
	}

	err = order.Assign(bestCourier.ID())
	if err != nil {
		return nil, err
	}

	return bestCourier, nil
}

func (c *CourierDispatcher) selectBestCourier(order *aggOrder.Order, couriers []*aggCourier.Courier) (*aggCourier.Courier, error) {
	var bestCourier *aggCourier.Courier
	minTime := math.MaxFloat64

	for _, courier := range couriers {
		if !courier.CanTakeOrder(order) {
			continue
		}

		timeToLocation := courier.CalculateTimeToLocation(order.Location())

		if timeToLocation < minTime {
			minTime = timeToLocation
			bestCourier = courier
		}
	}

	if bestCourier == nil {
		return nil, errs.NewValueIsInvalidErrorWithCause("couriers", errors.New("can't find a courier to take the order"))
	}

	return bestCourier, nil
}

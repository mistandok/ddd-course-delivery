package crons

import (
	"context"
	"log"

	"delivery/internal/core/application/usecases/commands/move_couriers_and_complete_order"
	"delivery/internal/pkg/errs"

	"github.com/robfig/cron/v3"
)

var _ cron.Job = &MoveCouriersJob{}

type MoveCouriersJob struct {
	moveCouriersCommandHandler move_couriers_and_complete_order.MoveCouriersAndCompleteOrderHandler
}

func NewMoveCouriersJob(
	moveCouriersCommandHandler move_couriers_and_complete_order.MoveCouriersAndCompleteOrderHandler) (cron.Job, error) {
	if moveCouriersCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("moveCouriersCommandHandler")
	}

	return &MoveCouriersJob{
		moveCouriersCommandHandler: moveCouriersCommandHandler}, nil
}

func (j *MoveCouriersJob) Run() {
	ctx := context.Background()
	command := move_couriers_and_complete_order.NewMoveCouriersAndFinishOrderCommand()

	err := j.moveCouriersCommandHandler.Handle(ctx, command)
	if err != nil {
		log.Printf("MoveCouriersJob error: %v", err)
	}
}

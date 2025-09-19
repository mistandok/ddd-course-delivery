package crons

import (
	"context"
	"log"

	"delivery/internal/core/application/usecases/commands/assign_order"
	"delivery/internal/pkg/errs"

	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrdersJob{}

type AssignOrdersJob struct {
	assignOrderCommandHandler assign_order.AssignedOrderHandler
}

func NewAssignOrdersJob(
	assignOrderCommandHandler assign_order.AssignedOrderHandler) (cron.Job, error) {
	if assignOrderCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("assignOrderCommandHandler")
	}

	return &AssignOrdersJob{
		assignOrderCommandHandler: assignOrderCommandHandler}, nil
}

func (j *AssignOrdersJob) Run() {
	ctx := context.Background()
	command := assign_order.NewAssignedOrderCommand()

	err := j.assignOrderCommandHandler.Handle(ctx, command)
	if err != nil {
		log.Printf("AssignOrdersJob error: %v", err)
	}
}

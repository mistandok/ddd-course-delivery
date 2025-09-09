package commands

import (
	"delivery/internal/pkg/errs"
	"errors"
)

type CreateCourierCommand struct {
	name  string
	speed int64

	isValid bool
}

func NewCreateCourierCommand(name string, speed int64) (CreateCourierCommand, error) {
	if name == "" {
		return CreateCourierCommand{}, errs.NewValueIsInvalidErrorWithCause("name", errors.New("name is required"))
	}

	if speed <= 0 {
		return CreateCourierCommand{}, errs.NewValueIsInvalidErrorWithCause("speed", errors.New("speed must be greater than 0"))
	}

	return CreateCourierCommand{name: name, speed: speed, isValid: true}, nil
}

func (c CreateCourierCommand) CommandName() string {
	return "CreateCourierCommand"
}

func (c CreateCourierCommand) IsValid() bool {
	return c.isValid
}

func (c CreateCourierCommand) Name() string {
	return c.name
}

func (c CreateCourierCommand) Speed() int64 {
	return c.speed
}

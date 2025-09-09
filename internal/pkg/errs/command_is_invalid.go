package errs

import (
	"errors"
	"fmt"
)

var ErrCommandIsInvalid = errors.New("command is invalid")

type CommandIsInvalidError struct {
	CommandName string
	Cause       error
}

func NewCommandIsInvalidErrorWithCause(commandName string, cause error) *CommandIsInvalidError {
	return &CommandIsInvalidError{
		CommandName: commandName,
		Cause:       cause,
	}
}

func NewCommandIsInvalidError(commandName string) *CommandIsInvalidError {
	return &CommandIsInvalidError{
		CommandName: commandName,
	}
}

func (e *CommandIsInvalidError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: command is: %s (cause: %v)",
			ErrCommandIsInvalid, e.CommandName, e.Cause)
	}
	return fmt.Sprintf("%s: %s", ErrCommandIsInvalid, e.CommandName)
}

func (e *CommandIsInvalidError) Unwrap() error {
	return ErrCommandIsInvalid
}

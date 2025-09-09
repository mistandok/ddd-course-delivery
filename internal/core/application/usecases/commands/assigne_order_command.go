package commands

type AssignedOrderCommand struct {
	isValid bool
}

func NewAssignedOrderCommand() AssignedOrderCommand {
	return AssignedOrderCommand{
		isValid: true,
	}
}

func (c AssignedOrderCommand) CommandName() string {
	return "AssignedOrderCommand"
}

func (c AssignedOrderCommand) IsValid() bool {
	return c.isValid
}

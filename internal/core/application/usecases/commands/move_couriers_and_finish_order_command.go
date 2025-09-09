package commands

type MoveCouriersAndFinishOrderCommand struct {
	isValid bool
}

func NewMoveCouriersAndFinishOrderCommand() MoveCouriersAndFinishOrderCommand {
	return MoveCouriersAndFinishOrderCommand{}
}

func (c MoveCouriersAndFinishOrderCommand) CommandName() string {
	return "MoveCouriersAndFinishOrderCommand"
}

func (c MoveCouriersAndFinishOrderCommand) IsValid() bool {
	return c.isValid
}

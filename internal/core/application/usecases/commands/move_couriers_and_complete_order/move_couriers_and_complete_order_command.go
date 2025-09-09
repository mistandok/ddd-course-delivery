package move_couriers_and_complete_order

type MoveCouriersAndFinishOrderCommand struct {
	isValid bool
}

func NewMoveCouriersAndFinishOrderCommand() MoveCouriersAndFinishOrderCommand {
	return MoveCouriersAndFinishOrderCommand{
		isValid: true,
	}
}

func (c MoveCouriersAndFinishOrderCommand) CommandName() string {
	return "MoveCouriersAndFinishOrderCommand"
}

func (c MoveCouriersAndFinishOrderCommand) IsValid() bool {
	return c.isValid
}

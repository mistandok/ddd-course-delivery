package get_all_uncompleted_orders

type GetAllUncompletedOrdersQuery struct {
	isValid bool
}

func NewGetAllUncompletedOrdersQuery() GetAllUncompletedOrdersQuery {
	return GetAllUncompletedOrdersQuery{isValid: true}
}

func (q GetAllUncompletedOrdersQuery) QueryName() string {
	return "GetAllUncompletedOrdersQuery"
}

func (q GetAllUncompletedOrdersQuery) IsValid() bool {
	return q.isValid
}

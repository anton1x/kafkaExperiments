package entities

type NewOrder struct {
	CustomerId string `json:"customerId"`
}

type Order struct {
	NewOrder
	Id string `json:"id"`
}

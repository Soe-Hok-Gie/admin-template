package domain

type Order struct {
	OrderID string
	Amount  int64
	Status  string
}

type StatusOrder struct {
	OrderID string
	Status  string
}

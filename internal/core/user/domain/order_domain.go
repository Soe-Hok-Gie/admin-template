package domain

type Order struct {
	OrderID string
	UserID  int64
	Amount  int64
	Status  string
}

type StatusOrder struct {
	OrderID string
	Status  string
}

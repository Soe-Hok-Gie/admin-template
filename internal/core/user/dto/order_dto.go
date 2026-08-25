package dto

type OrderRequest struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
	ItemName    string `json:"item_name"`
}

type OrderResponse struct {
	OrderID    string `json:"order_id"`
	PaymentURL string `json:"payment_url"`
	Status     string `json:"status"`
}

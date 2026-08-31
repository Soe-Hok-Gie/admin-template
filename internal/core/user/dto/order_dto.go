package dto

type OrderRequest struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
	ItemName    string `json:"item_name"`
	UserID      int64  `json:"_"`
}

type OrderResponse struct {
	OrderID    string `json:"order_id"`
	PaymentURL string `json:"payment_url"`
	Status     string `json:"status"`
}

type WebhookNotification struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	GrossAmount       int64  `json:"gross_amount"`
}

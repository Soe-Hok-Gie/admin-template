package service

import (
	"admin-template/internal/core/user/domain"
	"admin-template/internal/core/user/dto"
	"admin-template/internal/core/user/repository"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type OrderServiceImp struct {
	orderRepository repository.OrderRepository
	ServerKey       string // API Server Key dari Dashboard Midtrans Sandbox
}

func (service *OrderServiceImp) CreateOrder(ctx context.Context, req dto.OrderRequest) (*dto.OrderResponse, error) {
	// 1. Bungkus data ke dalam struct domain.Order sesuai kebutuhan fungsi SaveOrder
	orderData := domain.Order{
		OrderID: req.OrderID,
		Amount:  req.GrossAmount,
		Status:  "pending",
	}

	//save data order ke db
	err := service.orderRepository.SaveOrder(ctx, orderData)
	if err != nil {
		return nil, err
	}

	// Siapkan data untuk dikirim (Request Payload) ke Midtrans Snap API (menyiapkan document midtrans)
	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": req.GrossAmount,
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	//hit api Sandbox Midtrans Snap
	url := "https://midtrans.com"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	//service.ServerKey & Authorization adalah kunci rahasia/tanda tangan digital toko kita.
	auth := base64.StdEncoding.EncodeToString([]byte(service.ServerKey + ""))
	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Content-Type", "application/json")

	//menerima jawaban dari midtrans
	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Baca URL pembayaran (redirect_url) dari response Midtrans
	var midtransRes map[string]string
	json.NewDecoder(res.Body).Decode(&midtransRes)

	//output ke pembeli
	return &dto.OrderResponse{
		OrderID:    req.OrderID,
		PaymentURL: midtransRes["redirect_url"], // URL ini yang diberikan ke pembeli
		Status:     "PENDING",
	}, nil
}

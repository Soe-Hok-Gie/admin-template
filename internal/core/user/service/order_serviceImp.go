package service

import (
	"admin-template/internal/core/user/domain"
	"admin-template/internal/core/user/dto"
	"admin-template/internal/core/user/repository"
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type OrderServiceImp struct {
	orderRepository repository.OrderRepository
	ServerKey       string
}

func NewOrderService(
	orderRepository repository.OrderRepository,
	serverKey string,
) OrderService {
	return &OrderServiceImp{
		orderRepository: orderRepository,
		ServerKey:       strings.TrimSpace(serverKey),
	}

}

func (service *OrderServiceImp) CreateOrder(ctx context.Context, req dto.OrderRequest) (*dto.OrderResponse, error) {
	// 1. Bungkus data ke dalam struct domain.Order sesuai kebutuhan fungsi SaveOrder
	orderData := domain.Order{
		OrderID: req.OrderID,
		UserID:  req.UserID,
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
	url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	log.Printf("DEBUG URL MIDTRANS: %s", url)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	//service.ServerKey & Authorization adalah kunci rahasia/tanda tangan digital toko kita.
	auth := base64.StdEncoding.EncodeToString([]byte(service.ServerKey + ":"))
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
	var midtransRes map[string]interface{}
	json.NewDecoder(res.Body).Decode(&midtransRes)

	// 3. AMBIL REDIRECT URL DENGAN TYPE ASSERTION KE STRING
	paymentURL, ok := midtransRes["redirect_url"].(string)
	if !ok || paymentURL == "" {
		return nil, fmt.Errorf(
			"redirect_url tidak ditemukan dari Midtrans: %+v",
			midtransRes,
		)
	}

	//output ke pembeli
	return &dto.OrderResponse{
		OrderID:    req.OrderID,
		PaymentURL: paymentURL,
		Status:     "PENDING",
	}, nil
}

func (service *OrderServiceImp) ProcessWebhook(notification dto.WebhookNotification) error {
	// Ambil Server Key Anda dari konfigurasi/env
	ServerKey := os.Getenv("KEY_MIDTRANS")

	//verifikasi payload midtrans
	payloadSignature := notification.OrderID + notification.StatusCode + notification.GrossAmount + ServerKey
	hash := sha512.Sum512([]byte(payloadSignature))
	expectedSignatured := hex.EncodeToString(hash[:])

	if notification.SignatureKey != expectedSignatured {
		return fmt.Errorf("invalid signature key, request untrusted. Got: %s, Expected: %s",
			notification.SignatureKey,
			expectedSignatured,
		)
	}

	//Ambil data order asli dari database untuk cek nominal
	originalOrder, err := service.orderRepository.GetByOrderID(ctx, notification.OrderID)
	if err != nil {
		return fmt.Errorf("order not found in database :%w", err)
	}

	//terjemahin status midtrans ke app
	var finalStatus string
	switch notification.TransactionStatus {
	case "settlement", "capture":
		finalStatus = "PAID"
	case "expire", "cancel":
		finalStatus = "EXPIRED"
	}

	if finalStatus != "" {
		inputStatus := domain.StatusOrder{
			OrderID: notification.OrderID,
			Status:  finalStatus,
		}
		return service.orderRepository.UpdateStatus(inputStatus)
	}
	return nil
}

func (service *OrderServiceImp) CheckStatus(OrderID string) (string, error) {
	return service.orderRepository.GetStatus(OrderID)

}

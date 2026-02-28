package services

import (
	"crypto/sha512"
	"errors"
	"fmt"

	"greens-co/backend/internal/config"
	"greens-co/backend/internal/repositories"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService struct {
	orderRepo *repositories.OrderRepository
	cfg       *config.Config
}

func NewPaymentService(orderRepo *repositories.OrderRepository, cfg *config.Config) *PaymentService {
	return &PaymentService{orderRepo: orderRepo, cfg: cfg}
}

func (s *PaymentService) Initiate(orderID string) (string, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return "", errors.New("order not found")
	}

	env := midtrans.Sandbox
	if s.cfg.MidtransProduction {
		env = midtrans.Production
	}

	client := snap.Client{}
	client.New(s.cfg.MidtransServerKey, env)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.ID,
			GrossAmt: order.TotalPrice,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.Name,
			Phone: order.Phone,
		},
	}

	resp, err := client.CreateTransaction(req)
	if err != nil {
		return "", fmt.Errorf("midtrans error: %v", err)
	}

	// Store midtrans order ID
	s.orderRepo.SetMidtransID(order.ID, order.ID)

	return resp.RedirectURL, nil
}

func (s *PaymentService) HandleNotification(notification map[string]interface{}) error {
	// Verify signature: SHA512(order_id + status_code + gross_amount + server_key)
	orderID, _ := notification["order_id"].(string)
	statusCode, _ := notification["status_code"].(string)
	grossAmount, _ := notification["gross_amount"].(string)
	signatureKey, _ := notification["signature_key"].(string)

	expected := fmt.Sprintf("%x", sha512.Sum512([]byte(orderID+statusCode+grossAmount+s.cfg.MidtransServerKey)))
	if signatureKey != expected {
		return errors.New("invalid signature")
	}

	transactionStatus, _ := notification["transaction_status"].(string)
	fraudStatus, _ := notification["fraud_status"].(string)

	var newStatus, newPaymentStatus string
	if transactionStatus == "capture" {
		if fraudStatus == "accept" {
			newStatus = "PAID"
			newPaymentStatus = "PAID"
		}
	} else if transactionStatus == "settlement" {
		newStatus = "PAID"
		newPaymentStatus = "PAID"
	} else if transactionStatus == "cancel" || transactionStatus == "deny" || transactionStatus == "expire" {
		newStatus = "CANCELLED"
		newPaymentStatus = "UNPAID"
	}

	if newStatus != "" {
		order, err := s.orderRepo.FindByID(orderID)
		if err != nil {
			return err
		}
		s.orderRepo.UpdateStatus(order.ID, newStatus)
		s.orderRepo.UpdatePaymentStatus(order.ID, newPaymentStatus)
	}

	return nil
}

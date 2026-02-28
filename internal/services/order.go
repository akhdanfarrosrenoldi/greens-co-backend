package services

import (
	"greens-co/backend/internal/dto"
	"greens-co/backend/internal/models"
	"greens-co/backend/internal/repositories"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo   *repositories.OrderRepository
	productRepo *repositories.ProductRepository
}

func NewOrderService(orderRepo *repositories.OrderRepository, productRepo *repositories.ProductRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, productRepo: productRepo}
}

func (s *OrderService) Create(c echo.Context, req *dto.CreateOrderRequest) (*models.Order, error) {
	userID, _ := c.Get("userID").(string)

	// Calculate total
	var total int64
	for _, item := range req.Items {
		total += item.Price * int64(item.Qty)
	}
	// Delivery fee
	if req.Type == "DELIVERY" {
		total += 10000
	}

	order := &models.Order{
		UserID:        userID,
		Status:        "PENDING",
		Type:          req.Type,
		TotalPrice:    total,
		Name:          req.Name,
		Phone:         req.Phone,
		Address:       req.Address,
		PickupTime:    req.PickupTime,
		Notes:         req.Notes,
		PaymentStatus: "UNPAID",
	}

	for _, item := range req.Items {
		oi := models.OrderItem{
			ProductID: item.ProductID,
			VariantID: item.VariantID,
			Qty:       item.Qty,
			Price:     item.Price,
			Notes:     item.Notes,
		}
		order.Items = append(order.Items, oi)
	}

	// Run in a transaction: create order + decrement stock
	err := s.orderRepo.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.orderRepo.Create(tx, order); err != nil {
			return err
		}
		for _, item := range req.Items {
			if err := s.productRepo.DecrementStock(tx, item.ProductID, item.Qty); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(order.ID)
}

func (s *OrderService) GetMyOrders(c echo.Context) ([]models.Order, error) {
	userID, _ := c.Get("userID").(string)
	return s.orderRepo.FindByUserID(userID)
}

func (s *OrderService) GetAll(status string, page, limit int) ([]models.Order, int64, error) {
	return s.orderRepo.FindAll(status, page, limit)
}

func (s *OrderService) UpdateStatus(id, status string) (*models.Order, error) {
	if err := s.orderRepo.UpdateStatus(id, status); err != nil {
		return nil, err
	}
	return s.orderRepo.FindByID(id)
}

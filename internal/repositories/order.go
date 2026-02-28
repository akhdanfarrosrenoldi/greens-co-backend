package repositories

import (
	"greens-co/backend/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) preload(q *gorm.DB) *gorm.DB {
	return q.Preload("Items.Product.Category").
		Preload("Items.Product.Variants").
		Preload("Items.Variant").
		Preload("User")
}

func (r *OrderRepository) FindByUserID(userID string) ([]models.Order, error) {
	var orders []models.Order
	err := r.preload(r.db).Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindAll(status string, page, limit int) ([]models.Order, int64, error) {
	q := r.db.Model(&models.Order{})
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	q.Count(&total)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	var orders []models.Order
	err := r.preload(q).Order("created_at DESC").
		Offset((page - 1) * limit).Limit(limit).Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) FindByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.preload(r.db).First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByMidtransID(mid string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("midtrans_id = ?", mid).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Create(tx *gorm.DB, order *models.Order) error {
	return tx.Create(order).Error
}

func (r *OrderRepository) UpdateStatus(id, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *OrderRepository) UpdatePaymentStatus(id, paymentStatus string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).
		Update("payment_status", paymentStatus).Error
}

func (r *OrderRepository) SetMidtransID(id, midID string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).
		Update("midtrans_id", midID).Error
}

func (r *OrderRepository) CountTotal() (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Count(&count).Error
	return count, err
}

func (r *OrderRepository) SumRevenue() (int64, error) {
	var revenue int64
	err := r.db.Model(&models.Order{}).
		Where("status != ?", "CANCELLED").
		Select("COALESCE(SUM(total_price), 0)").Scan(&revenue).Error
	return revenue, err
}

func (r *OrderRepository) CountCustomers() (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Distinct("user_id").Count(&count).Error
	return count, err
}

func (r *OrderRepository) DB() *gorm.DB {
	return r.db
}

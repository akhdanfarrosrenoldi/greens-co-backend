package repositories

import (
	"fmt"
	"strings"

	"greens-co/backend/internal/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

type ProductFilter struct {
	Cat       string
	Sort      string
	Min       int64
	Max       int64
	Available *bool
	Diet      string
	Search    string
	Page      int
	Limit     int
}

func (r *ProductRepository) FindAll(f ProductFilter) ([]models.Product, int64, error) {
	q := r.db.Model(&models.Product{}).Preload("Category").Preload("Variants")

	if f.Cat != "" {
		q = q.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", f.Cat)
	}
	if f.Search != "" {
		q = q.Where("products.name ILIKE ?", "%"+f.Search+"%")
	}
	if f.Min > 0 {
		q = q.Where("products.base_price >= ?", f.Min)
	}
	if f.Max > 0 {
		q = q.Where("products.base_price <= ?", f.Max)
	}
	if f.Available != nil {
		q = q.Where("products.is_available = ?", *f.Available)
	}

	switch f.Sort {
	case "newest":
		q = q.Order("products.created_at DESC")
	case "price_asc":
		q = q.Order("products.base_price ASC")
	case "price_desc":
		q = q.Order("products.base_price DESC")
	default: // popular
		q = q.Order("COALESCE(products.review_count, 0) DESC")
	}

	var total int64
	q.Count(&total)

	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 9
	}
	offset := (f.Page - 1) * f.Limit
	q = q.Offset(offset).Limit(f.Limit)

	var products []models.Product
	err := q.Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) FindBySlug(slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("Variants").
		Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByID(id string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("Variants").First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(product).Error
}

func (r *ProductRepository) Delete(id string) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}

func (r *ProductRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *ProductRepository) UpdateVariants(productID string, variants []models.ProductVariant) error {
	// Delete variants not in new list
	if len(variants) > 0 {
		ids := make([]string, 0)
		for _, v := range variants {
			if v.ID != "" {
				ids = append(ids, v.ID)
			}
		}
		if len(ids) > 0 {
			r.db.Where("product_id = ? AND id NOT IN ?", productID, ids).Delete(&models.ProductVariant{})
		} else {
			r.db.Where("product_id = ?", productID).Delete(&models.ProductVariant{})
		}
	} else {
		r.db.Where("product_id = ?", productID).Delete(&models.ProductVariant{})
		return nil
	}

	for i := range variants {
		variants[i].ProductID = productID
	}
	return r.db.Save(&variants).Error
}

// GenerateUniqueSlug creates a slug from name and ensures it's unique
func (r *ProductRepository) GenerateUniqueSlug(name string) string {
	base := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, " ", "-"), "_", "-"))
	slug := base
	for i := 2; ; i++ {
		var count int64
		r.db.Model(&models.Product{}).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
	return slug
}

func (r *ProductRepository) DecrementStock(tx *gorm.DB, productID string, qty int) error {
	return tx.Model(&models.Product{}).
		Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock - ?", qty)).Error
}

package repositories

import (
	"fmt"
	"strings"

	"greens-co/backend/internal/models"

	"gorm.io/gorm"
)

type BundleRepository struct {
	db *gorm.DB
}

func NewBundleRepository(db *gorm.DB) *BundleRepository {
	return &BundleRepository{db: db}
}

func (r *BundleRepository) FindAll() ([]models.Bundle, error) {
	var bundles []models.Bundle
	err := r.db.Preload("Items.Product.Category").Preload("Items.Product.Variants").Find(&bundles).Error
	return bundles, err
}

func (r *BundleRepository) FindByID(id string) (*models.Bundle, error) {
	var bundle models.Bundle
	err := r.db.Preload("Items.Product.Category").Preload("Items.Product.Variants").
		First(&bundle, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (r *BundleRepository) Create(bundle *models.Bundle) error {
	return r.db.Create(bundle).Error
}

func (r *BundleRepository) Update(bundle *models.Bundle) error {
	// Replace items
	r.db.Where("bundle_id = ?", bundle.ID).Delete(&models.BundleItem{})
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(bundle).Error
}

func (r *BundleRepository) Delete(id string) error {
	r.db.Where("bundle_id = ?", id).Delete(&models.BundleItem{})
	return r.db.Delete(&models.Bundle{}, "id = ?", id).Error
}

func (r *BundleRepository) GenerateUniqueSlug(name string) string {
	base := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	slug := base
	for i := 2; ; i++ {
		var count int64
		r.db.Model(&models.Bundle{}).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
	return slug
}

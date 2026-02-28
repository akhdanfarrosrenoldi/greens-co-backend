package services

import (
	"greens-co/backend/internal/dto"
	"greens-co/backend/internal/models"
	"greens-co/backend/internal/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(f repositories.ProductFilter) ([]models.Product, int64, error) {
	return s.repo.FindAll(f)
}

func (s *ProductService) GetBySlug(slug string) (*models.Product, error) {
	return s.repo.FindBySlug(slug)
}

func (s *ProductService) Create(req *dto.CreateProductRequest) (*models.Product, error) {
	slug := req.Slug
	if slug == "" {
		slug = s.repo.GenerateUniqueSlug(req.Name)
	}

	product := &models.Product{
		Name:          req.Name,
		Slug:          slug,
		Description:   req.Description,
		BasePrice:     req.BasePrice,
		Stock:         req.Stock,
		Image:         req.Image,
		CategoryID:    req.CategoryID,
		IsAvailable:   req.IsAvailable,
		Badge:         req.Badge,
		OriginalPrice: req.OriginalPrice,
	}

	for _, v := range req.Variants {
		product.Variants = append(product.Variants, models.ProductVariant{
			Name:            v.Name,
			AdditionalPrice: v.AdditionalPrice,
		})
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return s.repo.FindByID(product.ID)
}

func (s *ProductService) Update(id string, req *dto.CreateProductRequest) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	slug := req.Slug
	if slug == "" {
		slug = s.repo.GenerateUniqueSlug(req.Name)
	}

	product.Name = req.Name
	product.Slug = slug
	product.Description = req.Description
	product.BasePrice = req.BasePrice
	product.Stock = req.Stock
	product.Image = req.Image
	product.CategoryID = req.CategoryID
	product.IsAvailable = req.IsAvailable
	product.Badge = req.Badge
	product.OriginalPrice = req.OriginalPrice

	// Sync variants
	var variants []models.ProductVariant
	for _, v := range req.Variants {
		vv := models.ProductVariant{
			Name:            v.Name,
			AdditionalPrice: v.AdditionalPrice,
			ProductID:       product.ID,
		}
		if v.ID != nil {
			vv.ID = *v.ID
		}
		variants = append(variants, vv)
	}
	if err := s.repo.UpdateVariants(product.ID, variants); err != nil {
		return nil, err
	}

	product.Variants = variants
	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return s.repo.FindByID(product.ID)
}

func (s *ProductService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *ProductService) GetAllAdmin(page, limit int) ([]models.Product, int64, error) {
	f := repositories.ProductFilter{Page: page, Limit: limit}
	return s.repo.FindAll(f)
}

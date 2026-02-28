package services

import (
	"greens-co/backend/internal/dto"
	"greens-co/backend/internal/models"
	"greens-co/backend/internal/repositories"
)

type BundleService struct {
	repo *repositories.BundleRepository
}

func NewBundleService(repo *repositories.BundleRepository) *BundleService {
	return &BundleService{repo: repo}
}

func (s *BundleService) GetAll() ([]models.Bundle, error) {
	return s.repo.FindAll()
}

func (s *BundleService) Create(req *dto.CreateBundleRequest) (*models.Bundle, error) {
	slug := req.Slug
	if slug == "" {
		slug = s.repo.GenerateUniqueSlug(req.Name)
	}

	bundle := &models.Bundle{
		Name:          req.Name,
		Slug:          slug,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Image:         req.Image,
		IsPopular:     req.IsPopular,
	}
	for _, item := range req.Items {
		bundle.Items = append(bundle.Items, models.BundleItem{
			ProductID: item.ProductID,
			Qty:       item.Qty,
		})
	}

	if err := s.repo.Create(bundle); err != nil {
		return nil, err
	}
	return s.repo.FindByID(bundle.ID)
}

func (s *BundleService) Update(id string, req *dto.CreateBundleRequest) (*models.Bundle, error) {
	bundle, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	slug := req.Slug
	if slug == "" {
		slug = s.repo.GenerateUniqueSlug(req.Name)
	}

	bundle.Name = req.Name
	bundle.Slug = slug
	bundle.Description = req.Description
	bundle.Price = req.Price
	bundle.OriginalPrice = req.OriginalPrice
	bundle.Image = req.Image
	bundle.IsPopular = req.IsPopular

	bundle.Items = nil
	for _, item := range req.Items {
		bundle.Items = append(bundle.Items, models.BundleItem{
			BundleID:  bundle.ID,
			ProductID: item.ProductID,
			Qty:       item.Qty,
		})
	}

	if err := s.repo.Update(bundle); err != nil {
		return nil, err
	}
	return s.repo.FindByID(bundle.ID)
}

func (s *BundleService) Delete(id string) error {
	return s.repo.Delete(id)
}

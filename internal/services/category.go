package services

import (
	"greens-co/backend/internal/dto"
	"greens-co/backend/internal/models"
	"greens-co/backend/internal/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) Create(req *dto.CreateCategoryRequest) (*models.Category, error) {
	slug := req.Slug
	if slug == "" {
		slug = s.repo.GenerateUniqueSlug(req.Name)
	}
	cat := &models.Category{Name: req.Name, Slug: slug, Image: req.Image}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Update(id string, req *dto.CreateCategoryRequest) (*models.Category, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	slug := req.Slug
	if slug == "" {
		slug = s.repo.GenerateUniqueSlug(req.Name)
	}
	cat.Name = req.Name
	cat.Slug = slug
	cat.Image = req.Image
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(id string) error {
	return s.repo.Delete(id)
}

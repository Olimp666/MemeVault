package service

import (
	"fmt"

	"github.com/Olimp666/MemeVault/internal/models"
)

type ImageRepository interface {
	Add(data []byte, userID int64, tags []string) (*models.Image, error)
	GetByTag(tag string) ([]*models.Image, error)
}

type Service struct {
	repo ImageRepository
}

func NewService(repo ImageRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) UploadImage(data []byte, userID int64, tags []string) (*models.Image, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("image data is empty")
	}

	if len(tags) == 0 {
		return nil, fmt.Errorf("at least one tag is required")
	}

	image, err := s.repo.Add(data, userID, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	return image, nil
}

func (s *Service) ImagesByTag(tag string) ([]*models.Image, error) {
	if tag == "" {
		return nil, fmt.Errorf("tag is required")
	}

	images, err := s.repo.GetByTag(tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get images by tag: %w", err)
	}

	return images, nil
}

package service

import (
	"fmt"

	"github.com/Olimp666/MemeVault/internal/models"
)

type ImageRepository interface {
	Add(tgFileID string, userID int64, tags []string) error
	GetByTags(tags []string, userID int64) ([]*models.Image, error)
}

type Service struct {
	repo ImageRepository
}

func NewService(repo ImageRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) UploadImage(tgFileID string, userID int64, tags []string) error {
	if tgFileID == "" {
		return fmt.Errorf("tg_file_id is empty")
	}

	if len(tags) == 0 {
		return fmt.Errorf("at least one tag is required")
	}

	err := s.repo.Add(tgFileID, userID, tags)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	return nil
}

func (s *Service) ImagesByTags(tags []string, userID int64) ([]*models.Image, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("at least one tag is required")
	}

	images, err := s.repo.GetByTags(tags, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get images by tags: %w", err)
	}

	return images, nil
}

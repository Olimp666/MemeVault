package service

import (
	"fmt"

	"github.com/Olimp666/MemeVault/internal/models"
)

type ImageRepository interface {
	Add(tgFileID string, userID int64, fileType string, tags []string) (int64, error)
	AddTags(imageID int64, tags []string) error
	ImageByUserAndFileID(userID int64, tgFileID string) (*models.Image, error)
	ImagesByTags(tags []string, userID int64) ([]*models.Image, error)
	ImagesBySubsetOfTags(tags []string, userID int64) ([]*models.Image, error)
	ImagesByUser(userID int64) ([]*models.ImageWithTags, error)
}

type Service struct {
	repo ImageRepository
}

func NewService(repo ImageRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) UploadImage(tgFileID string, userID int64, fileType string, tags []string) error {
	if tgFileID == "" {
		return fmt.Errorf("tg_file_id is empty")
	}

	if fileType == "" {
		return fmt.Errorf("file_type is empty")
	}

	if fileType != models.FileTypePhoto && fileType != models.FileTypeVideo && fileType != models.FileTypeGif {
		return fmt.Errorf("invalid file_type: must be 'photo', 'video', or 'gif'")
	}

	if len(tags) == 0 {
		return fmt.Errorf("at least one tag is required")
	}

	existingImage, err := s.repo.ImageByUserAndFileID(userID, tgFileID)
	if err == nil && existingImage != nil {
		err = s.repo.AddTags(existingImage.ID, tags)
		if err != nil {
			return fmt.Errorf("failed to add tags to existing image: %w", err)
		}
		return nil
	}

	_, err = s.repo.Add(tgFileID, userID, fileType, tags)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	return nil
}

func (s *Service) ImagesByTags(tags []string, userID int64) (exactMatch []*models.Image, partialMatch []*models.Image, err error) {
	if len(tags) == 0 {
		return nil, nil, fmt.Errorf("at least one tag is required")
	}

	exactMatch, err = s.repo.ImagesByTags(tags, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get images by tags: %w", err)
	}

	partialMatch, err = s.repo.ImagesBySubsetOfTags(tags, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get images by subset of tags: %w", err)
	}

	return exactMatch, partialMatch, nil
}

func (s *Service) ImagesByUser(userID int64) ([]*models.ImageWithTags, error) {
	images, err := s.repo.ImagesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get images by user: %w", err)
	}

	return images, nil
}

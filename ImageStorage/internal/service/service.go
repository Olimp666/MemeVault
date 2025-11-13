package service

import (
	"fmt"

	"github.com/Olimp666/MemeVault/internal/models"
)

type ImageRepository interface {
	Add(tgFileID string, userID int64, fileType string, tags []string) (int64, error)
	AddTags(imageID int64, tags []string) error
	ImageByUserAndFileID(userID int64, tgFileID string) (*models.Image, error)
	ImagesByTags(tags []string, userID int64) ([]*models.ImageWithTags, error)
	ImagesBySubsetOfTags(tags []string, userID int64) ([]*models.ImageWithTags, error)
	ImagesByUser(userID int64, sortBy string) ([]*models.ImageWithTags, error)
	DeleteImage(userID int64, tgFileID string) error
	DeleteAllUserImages(userID int64) error
	ReplaceTags(userID int64, tgFileID string, newTags []string) error
	IncrementUsageCount(userID int64, tgFileID string) error
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

func (s *Service) ImagesByTags(tags []string, userID int64) (exactMatch []*models.ImageWithTags, partialMatch []*models.ImageWithTags, err error) {
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

func (s *Service) ImagesByUser(userID int64, sortBy string) ([]*models.ImageWithTags, error) {
	if sortBy == "" {
		sortBy = models.SortByCreatedAt
	}
	
	if sortBy != models.SortByUsageCount && sortBy != models.SortByCreatedAt {
		return nil, fmt.Errorf("invalid sort_by parameter: must be '%s' or '%s'", models.SortByUsageCount, models.SortByCreatedAt)
	}

	images, err := s.repo.ImagesByUser(userID, sortBy)
	if err != nil {
		return nil, fmt.Errorf("failed to get images by user: %w", err)
	}

	return images, nil
}

func (s *Service) DeleteImage(userID int64, tgFileID string) error {
	if tgFileID == "" {
		return fmt.Errorf("tg_file_id is empty")
	}

	err := s.repo.DeleteImage(userID, tgFileID)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

func (s *Service) DeleteAllUserImages(userID int64) error {
	err := s.repo.DeleteAllUserImages(userID)
	if err != nil {
		return fmt.Errorf("failed to delete all user images: %w", err)
	}

	return nil
}

func (s *Service) ReplaceTags(userID int64, tgFileID string, newTags []string) error {
	if userID == models.DefaultUserID {
		return fmt.Errorf("cannot replace tags for default user")
	}

	if tgFileID == "" {
		return fmt.Errorf("tg_file_id is empty")
	}

	if len(newTags) == 0 {
		return fmt.Errorf("at least one tag is required")
	}

	err := s.repo.ReplaceTags(userID, tgFileID, newTags)
	if err != nil {
		return fmt.Errorf("failed to replace tags: %w", err)
	}

	return nil
}

func (s *Service) GenerateDescription(imageData []byte) (string, error) {
	return "метод пока не реализован", nil
}

func (s *Service) IncrementUsageCount(userID int64, tgFileID string) error {
	if tgFileID == "" {
		return fmt.Errorf("tg_file_id is empty")
	}

	err := s.repo.IncrementUsageCount(userID, tgFileID)
	if err != nil {
		return fmt.Errorf("failed to increment usage count: %w", err)
	}

	return nil
}

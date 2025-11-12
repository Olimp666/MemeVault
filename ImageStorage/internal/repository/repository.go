package repository

import (
	"fmt"

	"github.com/Olimp666/MemeVault/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) Add(tgFileID string, userID int64, fileType string, tags []string) (int64, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("can't begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var imageID int64
	queryImg := `INSERT INTO images (user_id, tg_file_id, file_type) VALUES ($1, $2, $3) RETURNING id;`
	err = tx.QueryRow(queryImg, userID, tgFileID, fileType).Scan(&imageID)
	if err != nil {
		return 0, fmt.Errorf("can't insert image: %w", err)
	}

	if len(tags) > 0 {
		queryTag := `INSERT INTO tags (image_id, name) VALUES (:image_id, :name)`

		for _, tag := range tags {
			_, err = tx.NamedExec(queryTag, map[string]interface{}{
				"image_id": imageID,
				"name":     tag,
			})
			if err != nil {
				return 0, fmt.Errorf("can't insert tag '%s': %w", tag, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("can't commit transaction: %w", err)
	}

	return imageID, nil
}

func (s *Repository) AddTags(imageID int64, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	queryTag := `INSERT INTO tags (image_id, name) VALUES (:image_id, :name) ON CONFLICT DO NOTHING`

	for _, tag := range tags {
		_, err = tx.NamedExec(queryTag, map[string]interface{}{
			"image_id": imageID,
			"name":     tag,
		})
		if err != nil {
			return fmt.Errorf("can't insert tag '%s': %w", tag, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}

	return nil
}

func (s *Repository) ImageByUserAndFileID(userID int64, tgFileID string) (*models.Image, error) {
	query := `SELECT id, user_id, tg_file_id, file_type, created_at FROM images WHERE user_id = $1 AND tg_file_id = $2;`

	var image models.Image
	err := s.db.Get(&image, query, userID, tgFileID)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (s *Repository) ImagesByTags(tags []string, userID int64) ([]*models.ImageWithTags, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("at least one tag is required")
	}

	query := `
		SELECT i.id, i.user_id, i.tg_file_id, i.file_type, i.created_at 
		FROM images i
		JOIN tags t ON t.image_id = i.id
		WHERE t.name = ANY($1) AND (i.user_id = $2 OR i.user_id = 0)
		GROUP BY i.id, i.user_id, i.tg_file_id, i.file_type, i.created_at
		HAVING COUNT(DISTINCT t.name) = $3
		ORDER BY i.created_at DESC;`

	var images []*models.Image
	err := s.db.Select(&images, query, pq.Array(tags), userID, len(tags))
	if err != nil {
		return nil, fmt.Errorf("can't get images by tags: %w", err)
	}

	result := make([]*models.ImageWithTags, 0, len(images))
	for _, img := range images {
		tagsQuery := `SELECT name FROM tags WHERE image_id = $1 ORDER BY name;`
		var imgTags []string
		err := s.db.Select(&imgTags, tagsQuery, img.ID)
		if err != nil {
			return nil, fmt.Errorf("can't get tags for image %d: %w", img.ID, err)
		}

		result = append(result, &models.ImageWithTags{
			Image: *img,
			Tags:  imgTags,
		})
	}

	return result, nil
}

func (s *Repository) ImagesBySubsetOfTags(tags []string, userID int64) ([]*models.ImageWithTags, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("at least one tag is required")
	}

	query := `
		SELECT i.id, i.user_id, i.tg_file_id, i.file_type, i.created_at
		FROM images i
		JOIN tags t ON t.image_id = i.id
		WHERE t.name = ANY($1) AND (i.user_id = $2 OR i.user_id = 0)
		GROUP BY i.id, i.user_id, i.tg_file_id, i.file_type, i.created_at
		HAVING COUNT(DISTINCT t.name) < $3
		ORDER BY COUNT(DISTINCT t.name) DESC, i.created_at DESC;`

	var images []*models.Image
	err := s.db.Select(&images, query, pq.Array(tags), userID, len(tags))
	if err != nil {
		return nil, fmt.Errorf("can't get images by subset of tags: %w", err)
	}

	result := make([]*models.ImageWithTags, 0, len(images))
	for _, img := range images {
		tagsQuery := `SELECT name FROM tags WHERE image_id = $1 ORDER BY name;`
		var imgTags []string
		err := s.db.Select(&imgTags, tagsQuery, img.ID)
		if err != nil {
			return nil, fmt.Errorf("can't get tags for image %d: %w", img.ID, err)
		}

		result = append(result, &models.ImageWithTags{
			Image: *img,
			Tags:  imgTags,
		})
	}

	return result, nil
}

func (s *Repository) ImagesByUser(userID int64) ([]*models.ImageWithTags, error) {
	query := `SELECT id, user_id, tg_file_id, file_type, created_at FROM images WHERE user_id = $1 ORDER BY created_at DESC;`

	var images []*models.Image
	err := s.db.Select(&images, query, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get images by user: %w", err)
	}

	result := make([]*models.ImageWithTags, 0, len(images))
	for _, img := range images {
		tagsQuery := `SELECT name FROM tags WHERE image_id = $1 ORDER BY name;`
		var tags []string
		err := s.db.Select(&tags, tagsQuery, img.ID)
		if err != nil {
			return nil, fmt.Errorf("can't get tags for image %d: %w", img.ID, err)
		}

		result = append(result, &models.ImageWithTags{
			Image: *img,
			Tags:  tags,
		})
	}

	return result, nil
}

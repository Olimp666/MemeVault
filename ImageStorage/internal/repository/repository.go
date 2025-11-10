package repository

import (
	"fmt"

	"github.com/Olimp666/MemeVault/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) Add(data []byte, userID int64, tags []string) (*models.Image, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("can't begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var imageID int64
	queryImg := `INSERT INTO images (user_id, data) VALUES ($1, $2) RETURNING id;`
	err = tx.QueryRowx(queryImg, userID, data).Scan(&imageID)
	if err != nil {
		return nil, fmt.Errorf("can't insert image: %w", err)
	}

	if len(tags) > 0 {
		queryTag := `INSERT INTO tags (image_id, name) VALUES (:image_id, :name)`

		for _, tag := range tags {
			_, err = tx.NamedExec(queryTag, map[string]interface{}{
				"image_id": imageID,
				"name":     tag,
			})
			if err != nil {
				return nil, fmt.Errorf("can't insert tag '%s': %w", tag, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("can't commit transaction: %w", err)
	}

	return &models.Image{
		ID:     imageID,
		UserID: userID,
		Data:   data,
	}, nil
}

func (s *Repository) GetByTag(tag string) ([]*models.Image, error) {
	const query = `SELECT i.id, i.user_id, i.data FROM images i JOIN tags t ON t.image_id = i.id WHERE t.name = $1;`

	var images []*models.Image
	err := s.db.Select(&images, query, tag)
	if err != nil {
		return nil, fmt.Errorf("can't get images by tag: %w", err)
	}

	return images, nil
}

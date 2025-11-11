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

func (s *Repository) Add(tgFileID string, userID int64, tags []string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	queryImg := `INSERT INTO images (tg_file_id, user_id) VALUES ($1, $2);`
	_, err = tx.Exec(queryImg, tgFileID, userID)
	if err != nil {
		return fmt.Errorf("can't insert image: %w", err)
	}

	if len(tags) > 0 {
		queryTag := `INSERT INTO tags (tg_file_id, name) VALUES (:tg_file_id, :name)`

		for _, tag := range tags {
			_, err = tx.NamedExec(queryTag, map[string]interface{}{
				"tg_file_id": tgFileID,
				"name":       tag,
			})
			if err != nil {
				return fmt.Errorf("can't insert tag '%s': %w", tag, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}

	return nil
}

func (s *Repository) GetByTags(tags []string, userID int64) ([]*models.Image, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("at least one tag is required")
	}

	query := `
		SELECT i.tg_file_id, i.user_id, i.created_at 
		FROM images i
		JOIN tags t ON t.tg_file_id = i.tg_file_id
		WHERE t.name = ANY($1) AND (i.user_id = $2 OR i.user_id = 0)
		GROUP BY i.tg_file_id, i.user_id, i.created_at
		HAVING COUNT(DISTINCT t.name) = $3
		ORDER BY i.created_at DESC;`

	var images []*models.Image
	err := s.db.Select(&images, query, pq.Array(tags), userID, len(tags))
	if err != nil {
		return nil, fmt.Errorf("can't get images by tags: %w", err)
	}

	return images, nil
}

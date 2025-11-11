package models

import "time"

type Image struct {
	TgFileID  string    `db:"tg_file_id"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

type Tag struct {
	TgFileID string `db:"tg_file_id"`
	Name     string `db:"name"`
}

type TagsRequest struct {
	Tags []string `json:"tags"`
}

type GetImagesResponse struct {
	TgFileIDs []string `json:"tg_file_ids"`
}

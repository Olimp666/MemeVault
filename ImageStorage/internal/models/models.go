package models

import "time"

const (
	FileTypePhoto = "photo"
	FileTypeVideo = "video"
	FileTypeGif   = "gif"
)

type Image struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	TgFileID  string    `db:"tg_file_id"`
	FileType  string    `db:"file_type"`
	CreatedAt time.Time `db:"created_at"`
}

type ImageWithTags struct {
	Image
	Tags []string
}

type Tag struct {
	ImageID int64  `db:"image_id"`
	Name    string `db:"name"`
}

type TagsRequest struct {
	Tags []string `json:"tags"`
}

type ImageInfo struct {
	TgFileID string `json:"tg_file_id"`
	FileType string `json:"file_type"`
}

type ImageInfoWithTags struct {
	TgFileID string   `json:"tg_file_id"`
	FileType string   `json:"file_type"`
	Tags     []string `json:"tags"`
}

type GetImagesResponse struct {
	ExactMatch   []ImageInfo `json:"exact_match"`
	PartialMatch []ImageInfo `json:"partial_match"`
}

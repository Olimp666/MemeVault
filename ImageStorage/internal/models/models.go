package models

import "time"

const (
	FileTypePhoto = "photo"
	FileTypeVideo = "video"
	FileTypeGif   = "gif"

	DefaultUserID = 0

	SortByUsageCount = "usage_count"
	SortByCreatedAt  = "created_at"
)

type Image struct {
	ID         int64     `db:"id"`
	UserID     int64     `db:"user_id"`
	TgFileID   string    `db:"tg_file_id"`
	FileType   string    `db:"file_type"`
	UsageCount int       `db:"usage_count"`
	CreatedAt  time.Time `db:"created_at"`
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
	ExactMatch   []ImageInfoWithTags `json:"exact_match"`
	PartialMatch []ImageInfoWithTags `json:"partial_match"`
}

type TelegramFileResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		FilePath string `json:"file_path"`
	} `json:"result"`
}

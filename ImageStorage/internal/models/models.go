package models

type Image struct {
	ID     int64  `db:"id"`
	UserID int64  `db:"user_id"`
	Data   []byte `db:"data"`
}

type Tag struct {
	ID      int64  `db:"id"`
	ImageID int64  `db:"image_id"`
	Name    string `db:"name"`
}

type UploadRequest struct {
	UserID int64    `json:"user_id"`
	Tags   []string `json:"tags"`
}

type UploadResponse struct {
	ImageID int64 `json:"image_id"`
}

type GetImagesResponse struct {
	Images []string `json:"images"`
}

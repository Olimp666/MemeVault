package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Olimp666/MemeVault/internal/models"
)

type ImageService interface {
	UploadImage(tgFileID string, userID int64, fileType string, tags []string) error
	ImagesByTags(tags []string, userID int64, sortBy string) (exactMatch []*models.ImageWithTags, partialMatch []*models.ImageWithTags, err error)
	ImagesByUser(userID int64, sortBy string) ([]*models.ImageWithTags, error)
	DeleteImage(userID int64, tgFileID string) error
	DeleteAllUserImages(userID int64) error
	ReplaceTags(userID int64, tgFileID string, newTags []string) error
	GenerateDescription(imageData []byte) (string, error)
	IncrementUsageCount(userID int64, tgFileID string) error
	ImageByTgFileID(tgFileID string) ([]byte, error)
}

type Handler struct {
	service ImageService
}

func NewHandler(service ImageService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	tgFileID := r.URL.Query().Get("tg_file_id")
	if tgFileID == "" {
		http.Error(w, "tg_file_id query parameter is required", http.StatusBadRequest)

		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)

		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)

		return
	}

	fileType := r.URL.Query().Get("file_type")
	if fileType == "" {
		http.Error(w, "file_type query parameter is required", http.StatusBadRequest)

		return
	}

	var req models.TagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request: %v", err), http.StatusBadRequest)

		return
	}

	err = h.service.UploadImage(tgFileID, userID, fileType, req.Tags)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to upload image: %v", err), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ImagesByTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)

		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)

		return
	}

	sortBy := r.URL.Query().Get("sort_by")

	var req models.TagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request: %v", err), http.StatusBadRequest)

		return
	}

	exactMatch, partialMatch, err := h.service.ImagesByTags(req.Tags, userID, sortBy)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get images: %v", err), http.StatusInternalServerError)

		return
	}

	response := models.GetImagesResponse{
		ExactMatch:   make([]models.ImageInfoWithTags, 0, len(exactMatch)),
		PartialMatch: make([]models.ImageInfoWithTags, 0, len(partialMatch)),
	}

	for _, img := range exactMatch {
		response.ExactMatch = append(response.ExactMatch, models.ImageInfoWithTags{
			TgFileID: img.TgFileID,
			FileType: img.FileType,
			Tags:     img.Tags,
		})
	}

	for _, img := range partialMatch {
		response.PartialMatch = append(response.PartialMatch, models.ImageInfoWithTags{
			TgFileID: img.TgFileID,
			FileType: img.FileType,
			Tags:     img.Tags,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ImagesByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)

		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)

		return
	}

	sortBy := r.URL.Query().Get("sort_by")

	images, err := h.service.ImagesByUser(userID, sortBy)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get images: %v", err), http.StatusInternalServerError)

		return
	}

	imageInfos := make([]models.ImageInfoWithTags, 0, len(images))
	for _, img := range images {
		imageInfos = append(imageInfos, models.ImageInfoWithTags{
			TgFileID: img.TgFileID,
			FileType: img.FileType,
			Tags:     img.Tags,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"images": imageInfos,
	})
}

func (h *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	tgFileID := r.URL.Query().Get("tg_file_id")
	if tgFileID == "" {
		http.Error(w, "tg_file_id query parameter is required", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteImage(userID, tgFileID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete image: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Image deleted successfully"))
}

func (h *Handler) DeleteAllUserImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteAllUserImages(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete all user images: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All user images deleted successfully"))
}

func (h *Handler) ReplaceTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	tgFileID := r.URL.Query().Get("tg_file_id")
	if tgFileID == "" {
		http.Error(w, "tg_file_id query parameter is required", http.StatusBadRequest)
		return
	}

	var req models.TagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request: %v", err), http.StatusBadRequest)
		return
	}

	err = h.service.ReplaceTags(userID, tgFileID, req.Tags)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to replace tags: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tags replaced successfully"))
}

func (h *Handler) GenerateDescription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageData := make([]byte, 10<<20)
	n, err := file.Read(imageData)
	if err != nil && err.Error() != "EOF" {
		http.Error(w, "failed to read image", http.StatusInternalServerError)
		return
	}
	imageData = imageData[:n]

	description, err := h.service.GenerateDescription(imageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate description: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"description": description,
	})
}

func (h *Handler) IncrementUsageCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	tgFileID := r.URL.Query().Get("tg_file_id")
	if tgFileID == "" {
		http.Error(w, "tg_file_id query parameter is required", http.StatusBadRequest)
		return
	}

	err = h.service.IncrementUsageCount(userID, tgFileID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to increment usage count: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usage count incremented successfully"))
}

func (h *Handler) ImageByFileID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tgFileID := r.PathValue("tg_file_id")
	if tgFileID == "" {
		http.Error(w, "tg_file_id is required", http.StatusBadRequest)
		return
	}

	imageData, err := h.service.ImageByTgFileID(tgFileID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(imageData)
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Olimp666/MemeVault/internal/models"
)

type ImageService interface {
	UploadImage(tgFileID string, userID int64, tags []string) error
	ImagesByTags(tags []string, userID int64) ([]*models.Image, error)
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

	var req models.TagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request: %v", err), http.StatusBadRequest)

		return
	}

	err = h.service.UploadImage(tgFileID, userID, req.Tags)
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

	var req models.TagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request: %v", err), http.StatusBadRequest)

		return
	}

	images, err := h.service.ImagesByTags(req.Tags, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get images: %v", err), http.StatusInternalServerError)

		return
	}

	response := models.GetImagesResponse{
		TgFileIDs: make([]string, 0, len(images)),
	}

	for _, img := range images {
		response.TgFileIDs = append(response.TgFileIDs, img.TgFileID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

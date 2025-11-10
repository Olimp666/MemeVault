package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Olimp666/MemeVault/internal/models"
)

type ImageService interface {
	UploadImage(data []byte, userID int64, tags []string) (*models.Image, error)
	ImagesByTag(tag string) ([]*models.Image, error)
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

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusBadRequest)

		return
	}

	userIDStr := r.FormValue("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)

		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)

		return
	}

	tagsJSON := r.FormValue("tags")
	if tagsJSON == "" {
		http.Error(w, "tags are required", http.StatusBadRequest)

		return
	}

	var tags []string
	if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
		http.Error(w, fmt.Sprintf("invalid tags format: %v", err), http.StatusBadRequest)

		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get image file: %v", err), http.StatusBadRequest)

		return
	}

	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read image: %v", err), http.StatusInternalServerError)

		return
	}

	image, err := h.service.UploadImage(imageData, userID, tags)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to upload image: %v", err), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.UploadResponse{
		ImageID: image.ID,
	})
}

func (h *Handler) ImagesByTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	tag := r.URL.Query().Get("tag")
	if tag == "" {
		http.Error(w, "tag parameter is required", http.StatusBadRequest)

		return
	}

	images, err := h.service.ImagesByTag(tag)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get images: %v", err), http.StatusInternalServerError)

		return
	}

	response := models.GetImagesResponse{
		Images: make([]string, 0, len(images)),
	}

	for _, img := range images {
		encoded := base64.StdEncoding.EncodeToString(img.Data)
		response.Images = append(response.Images, encoded)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

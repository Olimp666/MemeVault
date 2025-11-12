package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <путь_к_видео>")
		os.Exit(1)
	}

	videoPath := os.Args[1]
	outputPath := "description.txt"

	file, err := os.Open(videoPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("video", videoPath)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}
	writer.Close()

	resp, err := http.Post(
		"http://127.0.0.1:5000/caption_video",
		writer.FormDataContentType(),
		&body,
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result struct {
		CaptionEn string `json:"caption_en"`
		CaptionRu string `json:"caption_ru"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		panic(err)
	}

	fmt.Println("Описание (англ):", result.CaptionEn)
	fmt.Println("Описание (рус):", result.CaptionRu)

	tags := captionToTags(result.CaptionRu)
	fmt.Println("Теги:", tags)

	tagsJSON, _ := json.Marshal(tags)
	if err := os.WriteFile(outputPath, tagsJSON, 0o644); err != nil {
		panic(err)
	}
	fmt.Println("Теги записаны в", outputPath)
}

func captionToTags(caption string) []string {
	caption = strings.ToLower(caption)
	re := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	clean := re.ReplaceAllString(caption, "")
	words := strings.Fields(clean)
	tagMap := make(map[string]struct{})
	for _, w := range words {
		tagMap[w] = struct{}{}
	}
	var tags []string
	for t := range tagMap {
		tags = append(tags, t)
	}
	return tags
}

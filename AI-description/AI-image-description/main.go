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
		fmt.Println("Использование: go run main.go <путь_к_картинке>")
		os.Exit(1)
	}

	imagePath := os.Args[1]
	outputPath := "description.txt"

	file, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", imagePath)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}
	writer.Close()

	resp, err := http.Post("http://127.0.0.1:5000/caption", writer.FormDataContentType(), &body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Парсим JSON
	var result struct {
		Caption string `json:"caption"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		panic(err)
	}

	fmt.Println("Описание:", result.Caption)

	// =============================
	// 🔹 Преобразуем описание в список тегов
	// =============================
	tags := captionToTags(result.Caption)
	fmt.Println("Теги:", tags)

	// Сохраняем в файл (JSON-массив)
	tagsJSON, _ := json.Marshal(tags)
	if err := os.WriteFile(outputPath, tagsJSON, 0o644); err != nil {
		panic(err)
	}
	fmt.Println("Теги записаны в", outputPath)
}

// captionToTags превращает строку описания в массив тегов
func captionToTags(caption string) []string {
	// 1. Приводим к нижнему регистру
	caption = strings.ToLower(caption)

	// 2. Убираем пунктуацию, сохраняя русские буквы и цифры
	re := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	clean := re.ReplaceAllString(caption, "")

	// 3. Разбиваем на слова
	words := strings.Fields(clean)

	// 4. Для MVP — уникальные слова
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

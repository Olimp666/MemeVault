package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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

	resp, err := http.Post("http://127.0.0.1:5000/video_ocr", writer.FormDataContentType(), &body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		panic(err)
	}

	fmt.Println("Распознанный текст из видео:")
	fmt.Println(result.Text)

	if err := os.WriteFile(outputPath, []byte(result.Text), 0o644); err != nil {
		panic(err)
	}
	fmt.Println("Текст записан в", outputPath)
}

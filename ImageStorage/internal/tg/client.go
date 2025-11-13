package tg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Olimp666/MemeVault/internal/models"
)

const (
	telegramAPIURL = "https://api.telegram.org"
)

type Client struct {
	botToken string
}

func NewClient(botToken string) *Client {
	return &Client{
		botToken: botToken,
	}
}

func (c *Client) FilePath(fileID string) (string, error) {
	url := fmt.Sprintf("%s/bot%s/getFile?file_id=%s", telegramAPIURL, c.botToken, fileID)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("telegram API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var fileResp models.TelegramFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !fileResp.Ok {
		return "", fmt.Errorf("telegram API returned ok=false")
	}

	return fileResp.Result.FilePath, nil
}

func (c *Client) FileContent(filePath string) ([]byte, error) {
	url := fmt.Sprintf("%s/file/bot%s/%s", telegramAPIURL, c.botToken, filePath)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return data, nil
}

func (c *Client) ImageByFileID(fileID string) ([]byte, error) {
	filePath, err := c.FilePath(fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file path: %w", err)
	}

	imageData, err := c.FileContent(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download file content: %w", err)
	}

	return imageData, nil
}

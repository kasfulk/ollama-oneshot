package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *Client) ModelExists(modelName string) (bool, error) {
	tagsURL := c.tagsURL()
	resp, err := c.HTTPClient.Get(tagsURL)
	if err != nil {
		return false, fmt.Errorf("ollama tags request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("ollama tags returned %d: %s", resp.StatusCode, string(b))
	}

	var tags TagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return false, fmt.Errorf("decode tags response: %w", err)
	}

	for _, m := range tags.Models {
		if m.Name == modelName {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) tagsURL() string {
	host := c.BaseURL
	host = strings.TrimSuffix(host, "/api/generate")
	host = strings.TrimSuffix(host, "/")
	return host + "/api/tags"
}

func (c *Client) Generate(req *GenerateRequest) (string, error) {
	req.Stream = false
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(b))
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.Response, nil
}

func (c *Client) GenerateStream(req *GenerateRequest, onChunk func(chunk string)) error {
	req.Stream = true
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(b))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk StreamChunk
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			continue
		}
		if chunk.Response != "" {
			onChunk(chunk.Response)
		}
		if chunk.Done {
			break
		}
	}

	return scanner.Err()
}
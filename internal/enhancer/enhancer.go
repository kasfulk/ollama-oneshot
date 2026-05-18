package enhancer

import (
	"github.com/kasjfulk/ollama-oneshot/internal/ollama"
)

type Client struct {
	apiURL    string
	model     string
	ollama    *ollama.Client
}

func NewClient(apiURL, model string) *Client {
	return &Client{
		apiURL: apiURL,
		model:  model,
		ollama: ollama.NewClient(apiURL),
	}
}

func (c *Client) Enhance(userPrompt string) (string, error) {
	req := &ollama.GenerateRequest{
		Model:  c.model,
		Prompt: buildEnhancementPrompt(userPrompt),
		System: EnhancementSystemPrompt,
		Stream: false,
	}

	result, err := c.ollama.Generate(req)
	if err != nil {
		return "", err
	}

	if result == "" {
		return userPrompt, nil
	}

	return result, nil
}
package ollama

type StreamChunk struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at,omitempty"`
}

type TagsResponse struct {
	Models []ModelEntry `json:"models"`
}

type ModelEntry struct {
	Name string `json:"name"`
}
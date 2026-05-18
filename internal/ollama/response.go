package ollama

import "strings"

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
	Name        string `json:"name"`
	RemoteModel string `json:"remote_model,omitempty"`
	RemoteHost  string `json:"remote_host,omitempty"`
}

func (m *ModelEntry) IsCloud() bool {
	return m.RemoteHost != "" || strings.HasSuffix(m.Name, ":cloud")
}

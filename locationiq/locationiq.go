package locationiq

// APIError is the API error
type APIError struct {
	Error string `json:"error,omitempty"`
}

// NewClient makes a new client with the provided API key
func NewClient(key string) *Client {
	return &Client{key}
}

package monstercat

import "net/http"

// Monstercat Client.
type Client struct {
	httpClient *http.Client
}

// NewClient returns a new monstercat client
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = &http.Client{}
	}
	return &Client{c}
}

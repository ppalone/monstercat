package monstercat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	baseURL = "https://player.monstercat.app/api"
)

// Monstercat Client.
type Client struct {
	httpClient *http.Client
}

// NewClient returns a new monstercat client.
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = &http.Client{}
	}
	return &Client{c}
}

// SearchCatalog returns catalog search results for the provided query and optional search options.
func (c *Client) SearchCatalog(ctx context.Context, q string, opts ...Option) (SearchCatalogResults, error) {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}
	return c.searchCatalog(ctx, q, options)
}

func (c *Client) searchCatalog(ctx context.Context, q string, opts *options) (SearchCatalogResults, error) {
	opts.search = strings.TrimSpace(q)

	params, err := buildParams(opts)
	if err != nil {
		return SearchCatalogResults{}, err
	}

	req, err := makeRequest(ctx, "catalog/browse", params)
	if err != nil {
		return SearchCatalogResults{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchCatalogResults{}, err
	}
	defer resp.Body.Close()

	apiResponse := new(searchCatalogAPIResponse)
	err = json.NewDecoder(resp.Body).Decode(apiResponse)
	if err != nil {
		return SearchCatalogResults{}, err
	}

	return apiResponse.toResults(c, opts), nil
}

func makeRequest(ctx context.Context, url string, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", baseURL, url), nil)
	if err != nil {
		return nil, err
	}

	p := req.URL.Query()
	for k, v := range params {
		p.Set(k, v)
	}
	req.URL.RawQuery = p.Encode()

	return req, nil
}

func buildParams(opts *options) (map[string]string, error) {
	err := opts.validate()
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	params["limit"] = strconv.Itoa(opts.limit)
	params["offset"] = strconv.Itoa(opts.offset)
	params["search"] = opts.search
	params["sort"] = opts.sort

	if len(opts.releaseType) != 0 {
		params["types"] = opts.releaseType.String()
	}

	return params, nil
}

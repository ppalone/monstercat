package monstercat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	baseURL = "https://player.monstercat.app/api"
	webURL  = "https://www.monstercat.com"
	cdxURL  = "https://cdx.monstercat.com"
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

// GetTrackStream
func (c *Client) GetTrackStream(ctx context.Context, track Track) (io.ReadCloser, error) {
	if len(track.ID) == 0 {
		return nil, fmt.Errorf("track id is empty for track")
	}

	if len(track.Release.ID) == 0 {
		return nil, fmt.Errorf("release id is empty for track")
	}

	req, err := makeRequest(ctx, fmt.Sprintf("release/%s/track-stream/%s", track.Release.ID, track.ID), nil)
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		w.CloseWithError(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		errInvalidTrack := fmt.Errorf("invalid track")
		w.CloseWithError(errInvalidTrack)
		return nil, errInvalidTrack
	}

	go func() {
		defer resp.Body.Close()
		_, copyErr := io.Copy(w, resp.Body)
		w.CloseWithError(copyErr)
	}()

	return r, nil
}

// GetTrackStreamURL
func (c *Client) GetTrackStreamURL(ctx context.Context, track Track) (string, error) {
	if len(track.ID) == 0 {
		return "", fmt.Errorf("track id is empty for track")
	}

	if len(track.Release.ID) == 0 {
		return "", fmt.Errorf("release id is empty for track")
	}

	params := make(map[string]string)
	params["noRedirect"] = "true"

	req, err := makeRequest(ctx, fmt.Sprintf("release/%s/track-stream/%s", track.Release.ID, track.ID), params)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid track")
	}

	type trackStreamURL struct {
		SignedURL string `json:"SignedURL"`
	}
	res := new(trackStreamURL)
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return "", err
	}

	return res.SignedURL, nil
}

// GetResizedImage
func (c *Client) GetResizedImageURL(ctx context.Context, coverURL string, options ...ResizeOption) (string, error) {
	opts := newResizeOptions()
	for _, option := range options {
		option(opts)
	}

	if err := opts.validate(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cdxURL, nil)
	if err != nil {
		return "", err
	}

	params := req.URL.Query()
	params.Set("url", coverURL)
	params.Set("width", strconv.Itoa(opts.width))
	params.Set("encoding", string(opts.encoding))
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPermanentRedirect {
		return "", fmt.Errorf("failed to get resized image url")
	}

	location, err := resp.Location()
	if err != nil {
		return "", fmt.Errorf("error getting Location header: %w", err)
	}

	return location.String(), nil
}

// GetRelease
func (c *Client) GetRelease(ctx context.Context, id string, options ...ReleaseOption) (ReleaseInfo, error) {
	opts := newGetReleaseOpts()
	for _, option := range options {
		option(opts)
	}
	return c.getRelease(ctx, id, opts)
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

func (c *Client) getRelease(ctx context.Context, id string, opts *getReleaseOpts) (ReleaseInfo, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return ReleaseInfo{}, fmt.Errorf("id cannot be empty")
	}

	req, err := makeRequest(ctx, fmt.Sprintf("catalog/release/%s", id), opts.build())
	if err != nil {
		return ReleaseInfo{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ReleaseInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ReleaseInfo{}, fmt.Errorf("invalid id")
	}

	apiResponse := new(getReleaseAPIResponse)
	err = json.NewDecoder(resp.Body).Decode(apiResponse)
	if err != nil {
		return ReleaseInfo{}, err
	}

	return apiResponse.toReleaseInfo(ctx, c)
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

	if len(opts.releaseId) != 0 {
		params["releaseId"] = opts.releaseId
	}

	return params, nil
}

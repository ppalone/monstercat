package monstercat

import (
	"context"
	"fmt"
	"time"
)

type ReleaseType string

// release types.
const (
	ReleaseSingle ReleaseType = "Single"
	ReleaseEP     ReleaseType = "EP"
	ReleaseAlbum  ReleaseType = "Album"
)

func (t ReleaseType) String() string {
	return string(t)
}

// Release.
type Release struct {
	CatalogID string
	ID        string
	Title     string
	Type      string
	CoverURL  string
}

// Release API Response.
type releaseAPIResponse struct {
	ArtistsTitle        string    `json:"ArtistsTitle"`
	CatalogID           string    `json:"CatalogId"`
	Description         string    `json:"Description"`
	ID                  string    `json:"Id"`
	ReleaseDate         time.Time `json:"ReleaseDate"`
	ReleaseDateTimezone string    `json:"ReleaseDateTimezone"`
	Title               string    `json:"Title"`
	Type                string    `json:"Type"`
	UPC                 string    `json:"UPC"`
	Version             string    `json:"Version"`
}

// Release Info.
type ReleaseInfo struct {
	CatalogID string
	ID        string
	Title     string
	Type      string
	CoverURL  string
	Tracks    []Track
}

// Get Release API Response.
type getReleaseAPIResponse struct {
	Release struct {
		CacheTime            time.Time `json:"CacheTime"`
		CacheStatus          string    `json:"CacheStatus"`
		CacheStatusDetail    string    `json:"CacheStatusDetail"`
		AlbumNotes           string    `json:"AlbumNotes"`
		ArtistsTitle         string    `json:"ArtistsTitle"`
		BrandID              int       `json:"BrandId"`
		BrandTitle           string    `json:"BrandTitle"`
		CatalogID            string    `json:"CatalogId"`
		CopyrightPLine       string    `json:"CopyrightPLine"`
		CoverFileID          string    `json:"CoverFileId"`
		Description          string    `json:"Description"`
		FeaturedArtistsTitle string    `json:"FeaturedArtistsTitle"`
		GRid                 string    `json:"GRid"`
		GenrePrimary         string    `json:"GenrePrimary"`
		GenreSecondary       string    `json:"GenreSecondary"`
		ID                   string    `json:"Id"`
		PrereleaseDate       any       `json:"PrereleaseDate"`
		PresaveDate          any       `json:"PresaveDate"`
		ReleaseDate          time.Time `json:"ReleaseDate"`
		ReleaseDateTimezone  string    `json:"ReleaseDateTimezone"`
		SpotifyID            any       `json:"SpotifyId"`
		Title                string    `json:"Title"`
		Version              string    `json:"Version"`
		Type                 string    `json:"Type"`
		Upc                  string    `json:"UPC"`
		YouTubeURL           string    `json:"YouTubeUrl"`
		StreamingOnly        bool      `json:"StreamingOnly"`
		Freemium             bool      `json:"Freemium"`
		Downloadable         bool      `json:"Downloadable"`
		InEarlyAccess        bool      `json:"InEarlyAccess"`
		Streamable           bool      `json:"Streamable"`
	} `json:"Release"`
}

func (r *releaseAPIResponse) toRelease() Release {
	return Release{
		CatalogID: r.CatalogID,
		ID:        r.ID,
		Title:     r.Title,
		Type:      r.Type,
		CoverURL:  buildReleaseCoverURL(r.CatalogID),
	}
}

func (r *getReleaseAPIResponse) toReleaseInfo(ctx context.Context, c *Client) (ReleaseInfo, error) {
	tracks := make([]Track, 0)
	hasTracks := true

	for hasTracks {
		res, err := c.SearchCatalog(ctx, "", WithReleaseId(r.Release.ID))
		if err != nil {
			return ReleaseInfo{}, fmt.Errorf("error while getting tracks for release: %w", err)
		}

		tracks = append(tracks, res.Tracks...)
		hasTracks = res.HasNext
	}

	return ReleaseInfo{
		CatalogID: r.Release.CatalogID,
		ID:        r.Release.ID,
		Title:     r.Release.Title,
		Type:      r.Release.Type,
		CoverURL:  buildReleaseCoverURL(r.Release.CatalogID),
		Tracks:    tracks,
	}, nil
}

func buildReleaseCoverURL(catalogId string) string {
	return fmt.Sprintf("%s/release/%s/cover", webURL, catalogId)
}

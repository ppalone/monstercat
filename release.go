package monstercat

import (
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

func (r *releaseAPIResponse) toRelease() Release {
	return Release{
		CatalogID: r.CatalogID,
		ID:        r.ID,
		Title:     r.Title,
		Type:      r.Type,
		CoverURL:  fmt.Sprintf("%s/release/%s/cover", webURL, r.CatalogID),
	}
}

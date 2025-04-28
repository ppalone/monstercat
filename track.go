package monstercat

import (
	"fmt"
	"time"
)

// Track.
type Track struct {
	ID             string
	Title          string
	BrandID        int
	Brand          string
	DebutDate      time.Time
	BPM            int
	Duration       int
	Explicit       bool
	GenrePrimary   string
	GenreSecondary string
	Public         bool
	Release        Release
	ArtistsTitle   string
	Artists        []Artist
	StreamId       string // ReleaseCatalogID|TrackID
}

// Track API Response.
type trackAPIResponse struct {
	Artists         []artistAPIResponse `json:"Artists"`
	ArtistsTitle    string              `json:"ArtistsTitle"`
	BPM             int                 `json:"BPM"`
	Brand           string              `json:"Brand"`
	BrandID         int                 `json:"BrandId"`
	CreatorFriendly bool                `json:"CreatorFriendly"`
	DebutDate       time.Time           `json:"DebutDate"`
	Downloadable    bool                `json:"Downloadable"`
	Duration        int                 `json:"Duration"`
	Explicit        bool                `json:"Explicit"`
	Freemium        bool                `json:"Freemium"`
	GenrePrimary    string              `json:"GenrePrimary"`
	GenreSecondary  string              `json:"GenreSecondary"`
	ISRC            string              `json:"ISRC"`
	ID              string              `json:"Id"`
	InEarlyAccess   bool                `json:"InEarlyAccess"`
	LockStatus      string              `json:"LockStatus"`
	Public          bool                `json:"Public"`
	Release         releaseAPIResponse  `json:"Release"`
	Streamable      bool                `json:"Streamable"`
	StreamingOnly   bool                `json:"StreamingOnly"`
	Title           string              `json:"Title"`
	TrackNumber     int                 `json:"TrackNumber"`
	Version         string              `json:"Version"`
}

func (r *trackAPIResponse) toTrack() Track {
	artists := make([]Artist, 0)
	for _, result := range r.Artists {
		artists = append(artists, result.toArtist())
	}

	return Track{
		ID:             r.ID,
		Title:          r.Title,
		BrandID:        r.BrandID,
		Brand:          r.Brand,
		DebutDate:      r.DebutDate,
		BPM:            r.BPM,
		Duration:       r.Duration,
		Explicit:       r.Explicit,
		GenrePrimary:   r.GenrePrimary,
		GenreSecondary: r.GenreSecondary,
		Public:         r.Public,
		Release:        r.Release.toRelease(),
		ArtistsTitle:   r.ArtistsTitle,
		Artists:        artists,
		StreamId:       fmt.Sprintf("%s|%s", r.Release.ID, r.ID),
	}
}

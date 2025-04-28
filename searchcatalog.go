package monstercat

import (
	"context"
	"fmt"
)

// Search Catalog Results
type SearchCatalogResults struct {
	Limit   int
	Offset  int
	Size    int
	Total   int
	Tracks  []Track
	HasNext bool

	// for next
	c    *Client
	opts *options
}

// Search Catalog API Response
type searchCatalogAPIResponse struct {
	Limit  int                `json:"Limit"`
	Offset int                `json:"Offset"`
	Total  int                `json:"Total"`
	Data   []trackAPIResponse `json:"Data"`
}

func (r *searchCatalogAPIResponse) toResults(c *Client, opts *options) SearchCatalogResults {
	tracks := make([]Track, 0)
	for _, result := range r.Data {
		tracks = append(tracks, result.toTrack())
	}

	hasNext := (len(tracks) + r.Offset) < r.Total
	if !hasNext {
		return SearchCatalogResults{
			Limit:   r.Limit,
			Offset:  r.Offset,
			Size:    len(tracks),
			Total:   r.Total,
			Tracks:  tracks,
			HasNext: hasNext,
		}
	}

	return SearchCatalogResults{
		Limit:   r.Limit,
		Offset:  r.Offset,
		Size:    len(tracks),
		Total:   r.Total,
		Tracks:  tracks,
		HasNext: hasNext,
		c:       c,
		opts:    opts,
	}
}

func (results *SearchCatalogResults) Next(ctx context.Context) (SearchCatalogResults, error) {
	if !results.HasNext {
		return SearchCatalogResults{}, fmt.Errorf("no further results")
	}

	results.opts.offset += results.opts.limit
	return results.c.searchCatalog(ctx, results.opts.search, results.opts)
}

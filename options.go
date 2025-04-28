package monstercat

import "fmt"

type options struct {
	limit       int
	offset      int
	search      string
	sort        string
	releaseType ReleaseType
}

type Option func(o *options)

func newOptions() *options {
	return &options{
		limit:       100, // default limit // max 100
		offset:      0,
		search:      "",
		sort:        "",
		releaseType: "",
	}
}

func (o *options) validate() error {
	if o.limit < 1 || o.limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}

	return nil
}

func WithLimit(limit int) Option {
	return func(o *options) {
		o.limit = limit
	}
}

func WithOffset(offset int) Option {
	return func(o *options) {
		o.offset = offset
	}
}

func WithSort(field string) Option {
	return func(o *options) {
		o.sort = field
	}
}

func WithReleaseType(releaseType ReleaseType) Option {
	return func(o *options) {
		o.releaseType = releaseType
	}
}

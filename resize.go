package monstercat

import (
	"fmt"
)

type ImageEncoding string

const (
	JPEG ImageEncoding = "jpeg"
	WEBP ImageEncoding = "webp"
)

type resizeOptions struct {
	width    int
	encoding ImageEncoding
}

type ResizeOption func(o *resizeOptions)

func newResizeOptions() *resizeOptions {
	return &resizeOptions{
		width:    300,
		encoding: "webp",
	}
}

func (o *resizeOptions) validate() error {
	if !isEncodingAllowed(o.encoding) {
		return fmt.Errorf("invalid encoding")
	}

	return nil
}

func isEncodingAllowed(encoding ImageEncoding) bool {
	allowedEncodings := []ImageEncoding{JPEG, WEBP}
	for _, e := range allowedEncodings {
		if string(e) == string(encoding) {
			return true
		}
	}
	return false
}

func WithWidth(w int) ResizeOption {
	return func(o *resizeOptions) {
		o.width = w
	}
}

func WithEncoding(e ImageEncoding) ResizeOption {
	return func(o *resizeOptions) {
		o.encoding = e
	}
}

package monstercat_test

import (
	"context"
	"testing"

	"github.com/ppalone/monstercat"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	c := monstercat.NewClient(nil)
	assert.NotNil(t, c)
}

func Test_SearchCatalog(t *testing.T) {
	t.Run("with no options", func(t *testing.T) {
		q := "Nitro Fun"
		c := monstercat.NewClient(nil)
		res, err := c.SearchCatalog(context.Background(), q)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)

		for _, track := range res.Tracks {
			assert.NotEmpty(t, track.Artists)
		}
	})

	t.Run("with limit option", func(t *testing.T) {
		q := "Nitro Fun"
		c := monstercat.NewClient(nil)
		limit := 50
		res, err := c.SearchCatalog(context.Background(), q, monstercat.WithLimit(limit))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)
		assert.Equal(t, limit, res.Size)
	})

	t.Run("with offset option", func(t *testing.T) {
		q := ""
		c := monstercat.NewClient(nil)
		offset := 100
		res, err := c.SearchCatalog(context.Background(), q, monstercat.WithOffset(offset))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)
		assert.Equal(t, offset, res.Offset)
	})

	t.Run("with release type option", func(t *testing.T) {
		q := "Nitro Fun"
		c := monstercat.NewClient(nil)
		res, err := c.SearchCatalog(context.Background(), q, monstercat.WithReleaseType(monstercat.ReleaseSingle))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)

		for _, track := range res.Tracks {
			assert.Equal(t, monstercat.ReleaseSingle.String(), track.Release.Type)
		}
	})

	t.Run("with next and offset options", func(t *testing.T) {
		q := "Nitro Fun"
		c := monstercat.NewClient(nil)
		limit, offset := 50, 50
		res, err := c.SearchCatalog(context.Background(), q, monstercat.WithLimit(limit))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)
		assert.True(t, res.HasNext)

		resNext, err := res.Next(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, resNext.Tracks)

		resOffset, err := c.SearchCatalog(context.Background(), q, monstercat.WithLimit(limit), monstercat.WithOffset(offset))
		assert.NoError(t, err)
		assert.NotEmpty(t, resOffset.Tracks)

		assert.ElementsMatch(t, resNext.Tracks, resOffset.Tracks)
	})
}

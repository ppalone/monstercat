package monstercat_test

import (
	"context"
	"io"
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

	t.Run("with releaseId filter", func(t *testing.T) {
		c := monstercat.NewClient(nil)
		id := "475fcbbb-be8e-41bb-9f5e-1d3ce05f77be" // Best of 2024
		res, err := c.SearchCatalog(context.Background(), "", monstercat.WithReleaseId(id))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)
	})
}

func Test_GetTrackStream(t *testing.T) {
	t.Run("with track from catalog", func(t *testing.T) {
		c := monstercat.NewClient(nil)
		q := "Nitro Fun"
		res, err := c.SearchCatalog(context.Background(), q)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)

		track := res.Tracks[0]
		stream, err := c.GetTrackStream(context.Background(), track)
		assert.NoError(t, err)

		buf, err := io.ReadAll(stream)
		assert.NoError(t, err)
		assert.NotEmpty(t, buf)
	})

	t.Run("with custom track", func(t *testing.T) {
		// Pegboard Nerds - Emoji
		track := monstercat.Track{
			ID: "20330b9b-207f-47e9-9517-420fb0e32dbc",
			Release: monstercat.Release{
				ID: "81ae1308-60fe-4884-b3f4-23b8cf9d818e",
			},
		}

		c := monstercat.NewClient(nil)
		stream, err := c.GetTrackStream(context.Background(), track)
		assert.NoError(t, err)

		buf, err := io.ReadAll(stream)
		assert.NoError(t, err)
		assert.NotEmpty(t, buf)
	})

	t.Run("with invalid id", func(t *testing.T) {
		track := monstercat.Track{
			ID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			Release: monstercat.Release{
				ID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxx",
			},
		}

		c := monstercat.NewClient(nil)
		stream, err := c.GetTrackStream(context.Background(), track)
		assert.ErrorContains(t, err, "invalid track")
		assert.Nil(t, stream)
	})
}

func Test_GetTrackStreamURL(t *testing.T) {
	t.Run("with track from catalog", func(t *testing.T) {
		q := "Pegboard Nerds"
		c := monstercat.NewClient(nil)
		res, err := c.SearchCatalog(context.Background(), q, monstercat.WithReleaseType(monstercat.ReleaseSingle))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)

		track := res.Tracks[0]
		u, err := c.GetTrackStreamURL(context.Background(), track)
		assert.NoError(t, err)
		assert.NotEmpty(t, u)
	})

	t.Run("with invalid id", func(t *testing.T) {
		track := monstercat.Track{
			ID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			Release: monstercat.Release{
				ID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxx",
			},
		}

		c := monstercat.NewClient(nil)
		u, err := c.GetTrackStreamURL(context.Background(), track)
		assert.ErrorContains(t, err, "invalid track")
		assert.Empty(t, u)
	})
}

func Test_ReleaseCoverURL(t *testing.T) {
	c := monstercat.NewClient(nil)
	q := "Nitro Fun"
	res, err := c.SearchCatalog(context.Background(), q)
	assert.NoError(t, err)

	for _, track := range res.Tracks {
		assert.NotEmpty(t, track.Release.CoverURL)
	}
}

func Test_GetResizedImageURL(t *testing.T) {
	c := monstercat.NewClient(nil)
	q := "Nitro Fun"
	res, err := c.SearchCatalog(context.Background(), q, monstercat.WithLimit(5), monstercat.WithReleaseType(monstercat.ReleaseSingle))
	assert.NoError(t, err)

	for _, track := range res.Tracks {
		resizedURL, err := c.GetResizedImageURL(context.Background(), track.Release.CoverURL, monstercat.WithWidth(256))
		assert.NoError(t, err)
		assert.NotEmpty(t, resizedURL)
	}
}

func Test_GetRelease(t *testing.T) {
	t.Run("with catalog id", func(t *testing.T) {
		c := monstercat.NewClient(nil)
		id := "742779555328"
		res, err := c.GetRelease(context.Background(), id)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)
	})

	t.Run("with uuid", func(t *testing.T) {
		c := monstercat.NewClient(nil)
		uuid := "475fcbbb-be8e-41bb-9f5e-1d3ce05f77be"
		res, err := c.GetRelease(context.Background(), uuid, monstercat.WithIdType(monstercat.UUID))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Tracks)
	})

	t.Run("with invalid id", func(t *testing.T) {
		c := monstercat.NewClient(nil)
		id := "xxxxxxxx"
		_, err := c.GetRelease(context.Background(), id)
		assert.ErrorContains(t, err, "invalid id")
	})

	t.Run("with catalog id and uuid", func(t *testing.T) {
		c := monstercat.NewClient(nil)
		id := "742779555328"
		res1, err := c.GetRelease(context.Background(), id)
		assert.NoError(t, err)
		assert.NotEmpty(t, res1.Tracks)

		res2, err := c.GetRelease(context.Background(), res1.ID, monstercat.WithIdType(monstercat.UUID))
		assert.NoError(t, err)
		assert.NotEmpty(t, res2.Tracks)

		assert.Equal(t, res1.Title, res2.Title)
		assert.ElementsMatch(t, res1.Tracks, res2.Tracks)
	})
}

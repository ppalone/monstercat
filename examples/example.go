package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ppalone/monstercat"
)

func main() {
	c := monstercat.NewClient(nil)

	res, err := c.SearchCatalog(context.Background(), "Nitro Fun", monstercat.WithReleaseType(monstercat.ReleaseSingle))
	if err != nil {
		panic(err)
	}

	track := res.Tracks[0]
	stream, err := c.GetTrackStream(context.Background(), track)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fmt.Sprintf("./%s.mp3", track.Title))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := io.Copy(f, stream)
	if err != nil {
		panic(err)
	}

	fmt.Println("done:", n)
}

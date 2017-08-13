package main

import (
	"context"
	"fmt"
	"github.com/dyeduguru/wikiracer/racer"
	"github.com/dyeduguru/wikiracer/wikiclient/mediawiki"
	"time"
)

func main() {
	wikiClient := mediawiki.NewMediaWikiClient()
	wikiRacer := racer.NewGraphRacer(wikiClient)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	startTime := time.Now()
	path, err := wikiRacer.RaceWithTitle(ctx, "Mike Tyson", "Segment")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Time taken: %v\n", time.Since(startTime))
	fmt.Print("Path: ")
	for i, cur := range path {
		if i == len(path)-1 {
			fmt.Printf("%s", cur)
		} else {
			fmt.Printf("%s ->", cur)
		}
	}
}

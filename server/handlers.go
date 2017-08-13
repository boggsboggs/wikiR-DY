package server

import (
	"context"
	"encoding/json"
	"github.com/dyeduguru/wikiracer/racer"
	"github.com/dyeduguru/wikiracer/wikiclient/mediawiki"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	startVar = "start"
	endVar   = "end"
)

type response struct {
	path      []string
	timeTaken time.Duration
}

func RaceWithTitles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	start := vars[startVar]
	end := vars[endVar]
	wikiClient := mediawiki.NewMediaWikiClient()
	wikiRacer := racer.NewGraphRacer(wikiClient)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	startTime := time.Now()
	path, err := wikiRacer.RaceWithTitle(ctx, start, end)
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(response{
		path:      path,
		timeTaken: time.Since(startTime),
	}); err != nil {
		panic(err)
	}
}

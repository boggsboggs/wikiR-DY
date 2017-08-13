package server

import (
	"context"
	"encoding/json"
	"fmt"
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

type Response struct {
	Path      []string `json:"path"`
	TimeTaken float64  `json:"timeTaken"`
}

func RaceWithTitles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	start := vars[startVar]
	end := vars[endVar]
	fmt.Printf("Got request with start: %s and end: %s", start, end)
	wikiClient := mediawiki.NewMediaWikiClient()
	wikiRacer := racer.NewGraphRacer(wikiClient)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	startTime := time.Now()
	path, err := wikiRacer.RaceWithTitle(ctx, start, end)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&Response{
		Path:      path,
		TimeTaken: time.Since(startTime).Seconds(),
	}); err != nil {
		panic(err)
	}
}

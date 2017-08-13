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
	Error     string   `json:"error"`
}

func RaceWithTitles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	start := vars[startVar]
	end := vars[endVar]
	fmt.Printf("Got request with start: %s and end: %s\n", start, end)
	wikiClient := mediawiki.NewMediaWikiClient()
	wikiRacer := racer.NewGraphRacer(wikiClient)
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()
	startTime := time.Now()
	path, err := wikiRacer.RaceWithTitle(ctx, start, end)
	parsedResponse := &Response{}
	switch {
	case err == racer.NoPathError:
		parsedResponse.Error = "No Path between pages"
	case err == racer.TimedOutError:
		parsedResponse.Error = "Timed out"
	case err != nil:
		panic(err)
	}
	parsedResponse.Path = path
	parsedResponse.TimeTaken = time.Since(startTime).Seconds()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&parsedResponse); err != nil {
		panic(err)
	}
}

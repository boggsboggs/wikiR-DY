package server

import (
	"context"
	"encoding/base64"
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
	fmt.Printf("Got title request with start: %s and end: %s\n", start, end)
	racerFunc := func(ctx context.Context, r racer.Racer) ([]string, error) {
		return r.RaceWithTitle(ctx, start, end)
	}
	race(r.Context(), racerFunc, w)
}

func RaceWithURLs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	startByte, err := base64.StdEncoding.DecodeString(vars[startVar])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	endByte, err := base64.StdEncoding.DecodeString(vars[endVar])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	start, end := string(startByte), string(endByte)
	fmt.Printf("Got URL request with start: %s and end: %s\n", start, end)
	racerFunc := func(ctx context.Context, r racer.Racer) ([]string, error) {
		return r.RaceWithURL(ctx, start, end)
	}
	race(r.Context(), racerFunc, w)
}

func race(
	ctx context.Context,
	racerFunc func(ctx context.Context, r racer.Racer) ([]string, error),
	w http.ResponseWriter,
) {
	wikiClient := mediawiki.NewMediaWikiClient()
	wikiRacer := racer.NewGraphRacer(wikiClient)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	startTime := time.Now()
	path, err := racerFunc(ctx, wikiRacer)
	parsedResponse := &Response{}
	switch {
	case err == racer.NoPathError:
		parsedResponse.Error = "No Path between pages"
	case err == racer.TimedOutError:
		parsedResponse.Error = "Timed out"
	case err == racer.InvalidURLError:
		parsedResponse.Error = "Invalid URL"
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

package main

import (
	"fmt"
	"github.com/dyeduguru/wikiracer/graph"
	"github.com/dyeduguru/wikiracer/wikiclient/mediawiki"
)

func main() {
	wikiClient := mediawiki.NewMediaWikiClient()
	g := graph.New(wikiClient)
	path := g.Race("Palo Alto", "Nellore")
	for _, cur := range path {
		fmt.Printf("-> %s ", cur)
	}
}

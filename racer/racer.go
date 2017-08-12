package racer

import (
	"context"
	"errors"
	"github.com/dyeduguru/wikiracer/graph"
	"github.com/dyeduguru/wikiracer/wikiclient"
	"strings"
	"sync"
)

var (
	titlePrefixesToIgnore = map[string]struct{}{
		"Help:":                              {},
		"Wikipedia:":                         {},
		"International Standard Book Number": {},
		"Template":                           {},
		"Category":                           {},
	}
)

type Racer interface {
	RaceWithTitle(ctx context.Context, start, end string) ([]string, error)
	RaceWithURL(ctx context.Context, start, end string) ([]string, error)
}

type edge struct {
	src, dst string
}

func NewGraphRacer(client wikiclient.Client) Racer {
	return graphRacer{
		wikiClient:    client,
		graph:         graph.New(),
		leftCh:        make(chan edge),
		rightCh:       make(chan edge),
		leftFrontier:  make(map[*graph.Node]struct{}),
		rightFrontier: make(map[*graph.Node]struct{}),
	}
}

type graphRacer struct {
	wikiClient    wikiclient.Client
	graph         *graph.Graph
	m             sync.RWMutex
	rightFrontier map[*graph.Node]struct{}
	leftFrontier  map[*graph.Node]struct{}
	rightCh       chan edge
	leftCh        chan edge
}

func (g graphRacer) RaceWithTitle(ctx context.Context, start, end string) ([]string, error) {
	return g.race(ctx, start, end), nil
}

func (g graphRacer) RaceWithURL(ctx context.Context, start, end string) ([]string, error) {
	return nil, errors.New("Unimplemented")
}

func (g graphRacer) race(parentCtx context.Context, src, dst string) []string {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	if src == dst {
		return []string{src, dst}
	}
	go g.explore(ctx, src, g.leftCh)
	go g.explore(ctx, dst, g.rightCh)

	for {
		select {
		case e := <-g.leftCh:
			srcNode, dstNode := g.handleEdge(e, g.leftFrontier)
			if _, ok := g.rightFrontier[dstNode]; ok {
				return append(g.graph.Path(g.graph.LookUp[src], srcNode), reverse(g.graph.Path(g.graph.LookUp[dst], dstNode))...)
			}
			if dstNode.Title == dst {
				return g.graph.Path(g.graph.LookUp[src], dstNode)
			}
		case e := <-g.rightCh:
			srcNode, dstNode := g.handleEdge(e, g.rightFrontier)
			if _, ok := g.leftFrontier[dstNode]; ok {
				return append(g.graph.Path(g.graph.LookUp[src], dstNode), reverse(g.graph.Path(g.graph.LookUp[dst], srcNode))...)
			}
			if dstNode.Title == src {
				return reverse(g.graph.Path(g.graph.LookUp[dst], dstNode))
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (g graphRacer) handleEdge(e edge, frontier map[*graph.Node]struct{}) (*graph.Node, *graph.Node) {
	g.m.Lock()
	srcNode, dstNode := g.graph.InsertEdge(e.src, e.dst)
	g.m.Unlock()
	frontier[srcNode] = struct{}{}
	frontier[dstNode] = struct{}{}
	return srcNode, dstNode
}

func (g graphRacer) explore(ctx context.Context, start string, ch chan edge) {
	q := []string{start}
	visited := map[string]bool{}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if len(q) == 0 {
				return
			}
			cur := q[0]
			q = q[1:]
			visited[cur] = true
			titles, err := g.wikiClient.GetAllLinksInPage(cur)
			if err != nil {
				panic(err)
			}
			for _, title := range titles {
				if shouldIgnoreTitle(title) {
					continue
				}
				ch <- edge{src: cur, dst: title}
				if _, ok := visited[title]; ok {
					continue
				}
				q = append(q, title)
			}
		}
	}
}

func shouldIgnoreTitle(title string) bool {
	for titlePrefix := range titlePrefixesToIgnore {
		if strings.HasPrefix(title, titlePrefix) {
			return true
		}
	}
	return false
}

func reverse(s []string) []string {
	for i := 0; i <= len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
	return s
}

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
		"Wikipedia talk:":                    {},
		"International Standard Book Number": {},
		"Template":                           {},
		"Category":                           {},
		"User:":                              {},
		"User talk:":                         {},
		"Talk:":                              {},
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
		rightFinish:   make(chan struct{}),
		leftFinish:    make(chan struct{}),
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
	rightFinish   chan struct{}
	leftFinish    chan struct{}
}

func (g graphRacer) RaceWithTitle(ctx context.Context, start, end string) ([]string, error) {
	return g.race(ctx, start, end)
}

func (g graphRacer) RaceWithURL(ctx context.Context, start, end string) ([]string, error) {
	return nil, errors.New("Unimplemented")
}

func (g graphRacer) race(parentCtx context.Context, src, dst string) ([]string, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	if src == dst {
		return []string{src}, nil
	}
	go g.explore(ctx, src, g.leftCh, g.leftFinish, true)
	go g.explore(ctx, dst, g.rightCh, g.rightFinish, false)

	rightFinish, leftFinish := false, false
	for {
		select {
		case e := <-g.leftCh:
			srcNode, dstNode := g.handleEdge(e, g.leftFrontier)
			if _, ok := g.rightFrontier[dstNode]; ok {
				srcToLFrontier := g.graph.Path(g.graph.LookUp[src], srcNode)
				rFrontierToDst := g.graph.Path(dstNode, g.graph.LookUp[dst])
				return append(srcToLFrontier, rFrontierToDst...), nil
			}
			if dstNode.Title == dst {
				return g.graph.Path(g.graph.LookUp[src], dstNode), nil
			}
		case e := <-g.rightCh:
			srcNode, dstNode := g.handleEdge(e, g.rightFrontier)
			if _, ok := g.leftFrontier[dstNode]; ok {
				srcToLFrontier := g.graph.Path(g.graph.LookUp[src], dstNode)
				rFrontierToDst := g.graph.Path(srcNode, g.graph.LookUp[dst])
				return append(srcToLFrontier, rFrontierToDst...), nil
			}
			if dstNode.Title == src {
				return g.graph.Path(dstNode, g.graph.LookUp[dst]), nil
			}
		case <-g.leftFinish:
			leftFinish = true
		case <-g.rightFinish:
			rightFinish = true
		case <-ctx.Done():
			return nil, errors.New("timed out")
		}
		if rightFinish && leftFinish {
			return nil, errors.New("No path")
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

func (g graphRacer) explore(
	ctx context.Context,
	start string,
	ch chan edge,
	finishCh chan struct{},
	isForward bool,
) {
	q := []string{start}
	visited := map[string]bool{}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if len(q) == 0 {
				finishCh <- struct{}{}
				return
			}
			cur := q[0]
			q = q[1:]
			visited[cur] = true
			var titles []string
			var err error
			if isForward {
				titles, err = g.wikiClient.GetAllLinksFromPage(cur)
				if err != nil {
					panic(err)
				}
			} else {
				titles, err = g.wikiClient.GetAllLinksToPage(cur)
				if err != nil {
					panic(err)
				}
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

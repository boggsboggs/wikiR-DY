package graph

import (
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

type Node struct {
	Title     string
	Neighbors map[*Node]struct{}
}

type Edge struct {
	src, dst string
}

type Graph struct {
	m          sync.RWMutex
	LookUp     map[string]*Node
	FromRight  map[*Node]struct{}
	FromLeft   map[*Node]struct{}
	rightCh    chan *Edge
	leftCh     chan *Edge
	wikiClient wikiclient.Client
}

func New(client wikiclient.Client) *Graph {
	return &Graph{
		LookUp:     make(map[string]*Node),
		wikiClient: client,
		FromRight:  make(map[*Node]struct{}),
		FromLeft:   make(map[*Node]struct{}),
		leftCh:     make(chan *Edge),
		rightCh:    make(chan *Edge),
	}
}

func (g *Graph) Race(src, dst string) []string {
	if src == dst {
		return []string{src, dst}
	}
	go g.explore(src, g.leftCh)
	go g.explore(dst, g.rightCh)

	for {
		select {
		case edge := <-g.leftCh:
			srcNode, dstNode := g.insertEdge(edge)
			if _, ok := g.FromRight[dstNode]; ok {
				return append(g.path(g.LookUp[src], srcNode), reverse(g.path(g.LookUp[dst], dstNode))...)
			}
			if dstNode.Title == dst {
				return g.path(g.LookUp[src], dstNode)
			}
			g.FromLeft[srcNode] = struct{}{}
			g.FromLeft[dstNode] = struct{}{}
		case edge := <-g.rightCh:
			srcNode, dstNode := g.insertEdge(edge)
			if _, ok := g.FromLeft[dstNode]; ok {
				return append(g.path(g.LookUp[src], dstNode), reverse(g.path(g.LookUp[dst], srcNode))...)
			}

			if dstNode.Title == src {
				return reverse(g.path(g.LookUp[dst], dstNode))
			}
			g.FromRight[srcNode] = struct{}{}
			g.FromRight[dstNode] = struct{}{}
		}
	}
}

func (g *Graph) explore(start string, ch chan *Edge) {
	q := []string{start}
	for len(q) != 0 {
		cur := q[0]
		q = q[1:]
		titles, err := g.wikiClient.GetAllLinksInPage(cur)
		if err != nil {
			panic(err)
		}
		for _, title := range titles {
			if shouldIgnoreTitle(title) {
				continue
			}
			ch <- &Edge{src: cur, dst: title}
			if g.isTitlePresent(title) {
				continue
			}
			q = append(q, title)
		}
	}
}

func (g *Graph) path(src, dst *Node) []string {
	if src == dst {
		return []string{dst.Title}
	}
	q := []*Node{src}
	parent := map[*Node]*Node{}
	visited := map[*Node]struct{}{}
	for len(q) != 0 {
		cur := q[0]
		q = q[1:]
		visited[cur] = struct{}{}
		for neighbor := range cur.Neighbors {
			if _, ok := visited[neighbor]; ok {
				continue
			}
			parent[neighbor] = cur
			if neighbor == dst {
				path := []string{}
				cur = dst
				for cur != nil {
					path = append(path, cur.Title)
					cur = parent[cur]
				}
				return reverse(path)
			}
			q = append(q, neighbor)
		}
	}
	return nil
}

func (g *Graph) insertEdge(e *Edge) (*Node, *Node) {
	g.m.Lock()
	defer g.m.Unlock()
	if _, ok := g.LookUp[e.src]; !ok {
		g.LookUp[e.src] = newNode(e.src)
	}
	if _, ok := g.LookUp[e.dst]; !ok {
		g.LookUp[e.dst] = newNode(e.dst)
	}
	srcNode, dstNode := g.LookUp[e.src], g.LookUp[e.dst]
	srcNode.Neighbors[dstNode] = struct{}{}
	dstNode.Neighbors[srcNode] = struct{}{}
	return g.LookUp[e.src], g.LookUp[e.dst]
}

func (g *Graph) isTitlePresent(title string) bool {
	g.m.RLock()
	defer g.m.RUnlock()
	if _, ok := g.LookUp[title]; ok {
		return true
	}
	return false
}

func reverse(s []string) []string {
	ds := make([]string, len(s))
	for i := 0; i <= len(s)/2; i++ {
		ds[i], ds[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
	return ds
}

func newNode(title string) *Node {
	return &Node{
		Title:     title,
		Neighbors: make(map[*Node]struct{}),
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

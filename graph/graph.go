package graph

import "fmt"

type Node struct {
	Title     string
	Neighbors map[*Node]struct{}
}

type Graph struct {
	LookUp map[string]*Node
}

func New() *Graph {
	return &Graph{
		LookUp: make(map[string]*Node),
	}
}

func (g *Graph) Path(src, dst *Node) []string {
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
				reverse(path)
				return path
			}
			q = append(q, neighbor)
		}
	}
	return nil
}

func (g *Graph) InsertEdge(src, dst string) (*Node, *Node) {
	if _, ok := g.LookUp[src]; !ok {
		g.LookUp[src] = newNode(src)
	}
	if _, ok := g.LookUp[dst]; !ok {
		g.LookUp[dst] = newNode(dst)
	}
	srcNode, dstNode := g.LookUp[src], g.LookUp[dst]
	srcNode.Neighbors[dstNode] = struct{}{}
	dstNode.Neighbors[srcNode] = struct{}{}
	return g.LookUp[src], g.LookUp[dst]
}

func newNode(title string) *Node {
	return &Node{
		Title:     title,
		Neighbors: make(map[*Node]struct{}),
	}
}

func reverse(s []string) {
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
	fmt.Println()
}

package graph_test

import (
	"fmt"
	"github.com/dyeduguru/wikiracer/graph"
	"testing"
)

func TestGraph_Path(t *testing.T) {
	for _, tc := range []struct {
		initFunc     func(g *graph.Graph)
		src, dst     string
		expectedPath []string
	}{
		{
			initFunc: func(g *graph.Graph) {
				g.InsertEdge("foo", "bar")
				g.InsertEdge("bar", "foo")
				g.InsertEdge("bar", "baz")
			},
			src:          "foo",
			dst:          "baz",
			expectedPath: []string{"foo", "bar", "baz"},
		},
		{
			initFunc: func(g *graph.Graph) {
				g.InsertEdge("foo", "bar")
			},
			src:          "foo",
			dst:          "foo",
			expectedPath: []string{"foo"},
		},
		{
			initFunc: func(g *graph.Graph) {
				g.InsertEdge("foo", "bar")
			},
			src:          "foo",
			dst:          "baz",
			expectedPath: nil,
		},
	} {
		g := graph.New()
		tc.initFunc(g)
		path := g.Path(g.LookUp[tc.src], g.LookUp[tc.dst])
		if !testEq(path, tc.expectedPath) {
			fmt.Println("Actual path:")
			for _, val := range path {
				fmt.Printf("%s ", val)
			}
			fmt.Println()
			fmt.Println("Expected path:")
			for _, val := range tc.expectedPath {
				fmt.Printf("%s ", val)
			}
			fmt.Println()
			t.Errorf("returned path did not match excpected value")
		}
	}
}

func testEq(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

package racer_test

import (
	"context"
	"fmt"
	"github.com/dyeduguru/wikiracer/racer"
	"github.com/dyeduguru/wikiracer/wikiclient/mocks"
	"testing"
)

func defaultMockedClientFunc() *mocks.Client {
	mocksClient := &mocks.Client{}
	mocksClient.On("GetAllLinksFromPage", "foo").Return([]string{"bar"}, nil)
	mocksClient.On("GetAllLinksFromPage", "bar").Return([]string{"baz"}, nil)
	mocksClient.On("GetAllLinksFromPage", "baz").Return([]string{"qux"}, nil)
	mocksClient.On("GetAllLinksFromPage", "qux").Return([]string{}, nil)
	mocksClient.On("GetAllLinksToPage", "qux").Return([]string{"baz"}, nil)
	mocksClient.On("GetAllLinksToPage", "baz").Return([]string{"bar"}, nil)
	mocksClient.On("GetAllLinksToPage", "bar").Return([]string{"foo"}, nil)
	mocksClient.On("GetAllLinksToPage", "foo").Return([]string{}, nil)
	mocksClient.On("GetAllLinksToPage", "quxx").Return([]string{}, nil)
	return mocksClient
}

func TestGraphRacer_RaceWithTitle(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		src, dst        string
		getMockedClient func() *mocks.Client
		expectedPath    []string
	}{
		{
			getMockedClient: defaultMockedClientFunc,
			src:             "foo",
			dst:             "qux",
			expectedPath:    []string{"foo", "bar", "baz", "qux"},
		},
		{
			getMockedClient: defaultMockedClientFunc,
			src:             "foo",
			dst:             "bar",
			expectedPath:    []string{"foo", "bar"},
		},
		{
			getMockedClient: defaultMockedClientFunc,
			src:             "foo",
			dst:             "foo",
			expectedPath:    []string{"foo"},
		},
		{
			getMockedClient: defaultMockedClientFunc,
			src:             "foo",
			dst:             "quxx",
			expectedPath:    []string{},
		},
	} {
		racer := racer.NewGraphRacer(tc.getMockedClient())
		path, err := racer.RaceWithTitle(ctx, tc.src, tc.dst)
		if err != nil {
			t.Error(err.Error())
		}
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

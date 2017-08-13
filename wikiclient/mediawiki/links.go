package mediawiki

import (
	"encoding/json"
)

type LinksFromResponse struct {
	Continue struct {
		Plcontinue string `json:"plcontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Pages map[string]struct {
			Pageid int    `json:"pageid"`
			Ns     int    `json:"ns"`
			Title  string `json:"title"`
			Links  []struct {
				Ns    int    `json:"ns"`
				Title string `json:"title"`
			} `json:"links"`
		} `json:"pages"`
	} `json:"query"`
}

type LinksToResponse struct {
	Continue struct {
		Lhcontinue string `json:"lhcontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Pages map[string]struct {
			Pageid    int    `json:"pageid"`
			Ns        int    `json:"ns"`
			Title     string `json:"title"`
			LinksHere []struct {
				Ns    int    `json:"ns"`
				Title string `json:"title"`
			} `json:"linkshere"`
		} `json:"pages"`
	} `json:"query"`
}

type parsedLinkResponse struct {
	links         []string
	continueValue string
}

func mustParseLinksFromResponse(body []byte) parsedLinkResponse {
	resp := &LinksFromResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		panic(err)
	}
	links := []string{}
	for _, page := range resp.Query.Pages {
		for _, link := range page.Links {
			links = append(links, link.Title)
		}
	}
	return parsedLinkResponse{
		links:         links,
		continueValue: resp.Continue.Plcontinue,
	}
}

func mustParseLinksToResponse(body []byte) parsedLinkResponse {
	resp := &LinksToResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		panic(err)
	}
	links := []string{}
	for _, page := range resp.Query.Pages {
		for _, link := range page.LinksHere {
			links = append(links, link.Title)
		}
	}
	return parsedLinkResponse{
		links:         links,
		continueValue: resp.Continue.Lhcontinue,
	}
}

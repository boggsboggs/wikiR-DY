package mediawiki

type LinkResponse struct {
	Continue struct {
		Plcontinue string `json:"plcontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Normalized []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"normalized"`
		Pages map[string]Page `json:"pages"`
	} `json:"query"`
}

type Page struct {
	Pageid int    `json:"pageid"`
	Ns     int    `json:"ns"`
	Title  string `json:"title"`
	Links  []struct {
		Ns    int    `json:"ns"`
		Title string `json:"title"`
	} `json:"links"`
}

func getLinksFromResponse(resp *LinkResponse) []string {
	links := []string{}
	for _, page := range resp.Query.Pages {
		for _, link := range page.Links {
			links = append(links, link.Title)
		}
	}
	return links
}

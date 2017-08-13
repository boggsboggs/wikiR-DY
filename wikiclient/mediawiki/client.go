package mediawiki

import (
	"github.com/dyeduguru/wikiracer/wikiclient"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	baseURL = "https://en.wikipedia.org/w/api.php"
)

type mediaWikiClient struct {
	baseURL string
	client  *http.Client
}

func NewMediaWikiClient() wikiclient.Client {
	return mediaWikiClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Second * 2,
		},
	}
}

func (m mediaWikiClient) GetAllLinksFromPage(title string) ([]string, error) {
	params := map[string]string{
		"action":  "query",
		"titles":  title,
		"format":  "json",
		"prop":    "links",
		"pllimit": "500",
	}
	return m.getLinksWithContinue(params, "plcontinue", mustParseLinksFromResponse)
}

func (m mediaWikiClient) GetAllLinksToPage(title string) ([]string, error) {
	params := map[string]string{
		"action":  "query",
		"titles":  title,
		"format":  "json",
		"prop":    "linkshere",
		"lhlimit": "500",
	}
	return m.getLinksWithContinue(params, "lhcontinue", mustParseLinksToResponse)
}

func (m mediaWikiClient) getLinksWithContinue(
	params map[string]string,
	continuePrefix string,
	parseBodyFunc func([]byte) parsedLinkResponse,
) ([]string, error) {
	links := []string{}
	for {
		body, err := m.performAction(params)
		if err != nil {
			return nil, err
		}
		parsedBody := parseBodyFunc(body)
		links = append(links, parsedBody.links...)

		if parsedBody.continueValue == "" {
			break
		}
		params[continuePrefix] = parsedBody.continueValue
	}
	return links, nil
}

func (m mediaWikiClient) performAction(params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", m.baseURL, nil)
	if err != nil {
		return nil, err
	}
	values := req.URL.Query()
	for k, v := range params {
		values.Add(k, v)
	}
	req.URL.RawQuery = values.Encode()
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

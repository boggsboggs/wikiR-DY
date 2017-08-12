package mediawiki

import (
	"encoding/json"
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

func (m mediaWikiClient) GetAllLinksInPage(title string) ([]string, error) {
	params := map[string]string{
		"action":  "query",
		"titles":  title,
		"format":  "json",
		"prop":    "links",
		"pllimit": "5000",
	}

	links := []string{}
	for {
		body, err := m.performAction(params)
		if err != nil {
			return nil, err
		}
		resp := &LinkResponse{}
		if err := json.Unmarshal(body, resp); err != nil {
			return nil, err
		}
		links = append(links, getLinksFromResponse(resp)...)
		if resp.Continue.Plcontinue == "" {
			break
		}
		params["plcontinue"] = resp.Continue.Plcontinue
	}
	return links, nil
}

func (m mediaWikiClient) GetAllLinksInURL(url string) ([]string, error) {
	return nil, nil
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

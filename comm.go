package lastfm

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	apiBaseURL = "http://ws.audioscrobbler.com/2.0/?"
)

func buildQueryURL(query map[string]string) string {
	parts := make([]string, 0, len(query))
	for key, value := range query {
		parts = append(parts, strings.Join([]string{key, strings.Replace(value, " ", "+", -1)}, "="))
	}
	return apiBaseURL + strings.Join(parts, "&")
}

type mockServer interface {
	doQuery(params map[string]string) ([]byte, error)
}

type LastFM struct {
	apiKey string
	mock   mockServer
}

func New(apiKey string) LastFM {
	return LastFM{apiKey: apiKey}
}

func (lfm LastFM) doQuery(method string, params map[string]string) (data []byte, err error) {
	queryParams := make(map[string]string, len(params)+2)
	queryParams["api_key"] = lfm.apiKey
	queryParams["method"] = method
	for key, value := range params {
		queryParams[key] = value
	}
	if lfm.mock != nil {
		return lfm.mock.doQuery(queryParams)
	}

	resp, err := http.Get(buildQueryURL(queryParams))
	if err != nil {
		return
	}
	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

type lfmStatus struct {
	Status       string       `xml:"status,attr"`
	RecentTracks RecentTracks `xml:"recenttracks"`
	Tasteometer  Tasteometer  `xml:"comparison"`
	TrackInfo    TrackInfo    `xml:"track"`
	TopTags      TopTags      `xml:"toptags"`
	Neighbours   []Neighbour  `xml:"neighbours>user"`
	TopArtists   TopArtists   `xml:"topartists"`
}

type lfmDate struct {
	Date string `xml:",chardata"`
	UTS  int64  `xml:"uts,attr"`
}

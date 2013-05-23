package lastfm

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

type MockLastFM struct{}

func (_ *MockLastFM) Get(uri string) (resp *http.Response, err error) {
	args := strings.Split(strings.Replace(uri, apiBaseURL, "", 1), "&")
	query := make(map[string]string)
	for _, arg := range args {
		param := strings.Split(arg, "=")
		k, _ := url.QueryUnescape(param[0])
		query[k], _ = url.QueryUnescape(param[1])
	}
	fh, err := os.Open(buildMockFilename(query))

	if err != nil && os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Produce with: curl -o - '"+buildQueryURL(query)+"' > '"+buildMockFilename(query)+"'")
		return
	}
	resp = &http.Response{Body: fh} // doQuery cares about the Body
	return
}

func buildMockFilename(query map[string]string) string {
	parts := make([]string, 0, len(query)+1)
	parts = append(parts, query["method"])

	keys := make([]string, 0, len(query)-1)
	for key, _ := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if key == "method" || key == "api_key" {
			continue
		}
		parts = append(parts, strings.Join([]string{key, strings.Replace(query[key], " ", ".", -1)}, "="))
	}
	parts = append(parts, "xml")

	return "fixtures/" + strings.Join(parts, ".")
}

func Mock(lfm LastFM) LastFM {
	lfm.getter = &MockLastFM{}
	return lfm
}

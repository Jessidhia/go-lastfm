package lastfm

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type MockServer struct {
	mockServer
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

func (_ *MockServer) doQuery(params map[string]string) (data []byte, err error) {
	fh, err := os.Open(buildMockFilename(params))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Produce with: curl -o - '"+buildQueryURL(params)+"' > '"+buildMockFilename(params)+"'")
		}
		return
	}
	data, err = ioutil.ReadAll(fh)
	fh.Close()
	return
}

func Mock(lfm LastFM) LastFM {
	lfm.mock = &MockServer{}
	return lfm
}

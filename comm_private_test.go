package lastfm

import "testing"

func TestBuildQueryURL(T *testing.T) {
	T.Parallel()
	expect := "http://ws.audioscrobbler.com/2.0/?a=b&c=d"
	sample := buildQueryURL(map[string]string{
		"a": "b",
		"c": "d"})
	if sample != expect {
		T.Error("Expected sample query URL \"" + expect + "\" -- Got \"" + sample + "\"")
	}
	return
}

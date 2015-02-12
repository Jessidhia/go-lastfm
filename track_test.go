package lastfm

import (
	"testing"
	"time"
)

func TestGetTrackInfo_ByMBID(T *testing.T) {
	T.Parallel()
	lfm := testMock(New("4c563adf68bc357a4570d3e7986f6481"))
	trackInfo, err := lfm.GetTrackInfo(
		Track{MBID: "29b45fae-fc32-43c0-ab74-052842458315"}, "", false)

	if testExpect(T, "error", nil, err) {
		testExpect(T, "track ID", 4313, trackInfo.ID)
		testExpect(T, "artist", "Daft Punk", trackInfo.Artist.Name)
		testExpect(T, "album", "Discovery", trackInfo.Album.Name)
		testExpect(T, "top tag", "electronic", trackInfo.TopTags[0])
		dur, _ := time.ParseDuration("212s")
		testExpect(T, "duration", dur, trackInfo.Duration)
		testExpect(T, "user playcount", 0, trackInfo.UserPlaycount)
	}
}

func TestGetTrackInfo_ByTrackArtist(T *testing.T) {
	T.Parallel()
	lfm := testMock(New("4c563adf68bc357a4570d3e7986f6481"))
	trackInfo, err := lfm.GetTrackInfo(
		Track{Artist: Artist{Name: "Daft Punk"}, Name: "Motherboard"},
		"Kovensky", false)

	if testExpect(T, "error", nil, err) {
		testExpect(T, "track ID", 651481384, trackInfo.ID)
		testExpect(T, "album", "Random Access Memories", trackInfo.Album.Name)
		testExpect(T, "tag count", 0, len(trackInfo.TopTags))
		dur, _ := time.ParseDuration("326s")
		testExpect(T, "duration", dur, trackInfo.Duration)
		testExpect(T, "loved", true, trackInfo.UserLoved)
		testExpect(T, "user playcount", 64, trackInfo.UserPlaycount)
	}
}

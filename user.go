package lastfm

import (
	"encoding/xml"
	"strconv"
)

type RecentTracks struct {
	User       string  `xml:"user,attr"`
	Total      int     `xml:"total,attr"`
	Tracks     []Track `xml:"track"`
	NowPlaying *Track  `xml:"-"` // Points to the currently playing track, if any
}

func (tracks *RecentTracks) unmarshalHelper() (err error) {
	for i, track := range tracks.Tracks {
		if track.NowPlaying {
			tracks.NowPlaying = &tracks.Tracks[i]
		}
		err = tracks.Tracks[i].unmarshalHelper()
		if err != nil {
			return
		}
	}
	return
}

// Gets a list of recent tracks from the user. The .Tracks field includes the currently playing track,
// if any, and up to the count most recent scrobbles.
// The .NowPlaying field points to any currently playing track.
//
// See http://www.last.fm/api/show/user.getRecentTracks.
func (lfm LastFM) GetRecentTracks(user string, count int) (tracks *RecentTracks, err error) {
	bytes, err := lfm.doQuery("user.getRecentTracks", map[string]string{
		"user":     user,
		"extended": "1",
		"limit":    strconv.Itoa(count)})
	if err != nil {
		return
	}
	// Sadly, this errors out because <recenttracks> isn't the toplevel element:
	// tracks = &RecentTracks{}
	// err = xml.Unmarshal(bytes, tracks)
	// Using `xml:"recenttracks>track"` works for .Tracks, but not for .User or .Total
	status := lfmStatus{}
	err = xml.Unmarshal(bytes, &status)
	if err != nil {
		return
	}
	if status.Error.Code != 0 {
		err = &status.Error
		return
	}

	tracks = &status.RecentTracks
	err = tracks.unmarshalHelper()
	return
}

type Tasteometer struct {
	Users   []string `xml:"input>user>name"`            // The compared users
	Score   float32  `xml:"result>score"`               // Varies from 0.0 to 1.0
	Artists []string `xml:"result>artists>artist>name"` // Short list of up to 5 common artists with the most affinity
}

// Compares the taste of 2 users.
//
// See http://www.last.fm/api/show/tasteometer.compare.
func (lfm LastFM) CompareTaste(user1 string, user2 string) (taste *Tasteometer, err error) {
	bytes, err := lfm.doQuery("tasteometer.compare", map[string]string{
		"type1":  "user",
		"type2":  "user",
		"value1": user1,
		"value2": user2})
	if err != nil {
		return
	}
	status := lfmStatus{}
	err = xml.Unmarshal(bytes, &status)
	if err != nil {
		return
	}
	if status.Error.Code != 0 {
		err = &status.Error
		return
	}

	taste = &status.Tasteometer
	return
}

type Neighbour struct {
	Name  string  `xml:"name"`
	Match float32 `xml:"match"`
}

// Gets a list of up to limit closest neighbours of a user. A neighbour is another user
// that has high tasteometer comparison scores.
//
// See http://www.last.fm/api/show/user.getNeighbours
func (lfm LastFM) GetUserNeighbours(user string, limit int) (neighbours []Neighbour, err error) {
	bytes, err := lfm.doQuery("user.getNeighbours", map[string]string{
		"user":  user,
		"limit": strconv.Itoa(limit)})
	if err != nil {
		return
	}
	status := lfmStatus{}
	err = xml.Unmarshal(bytes, &status)
	if err != nil {
		return
	}
	if status.Error.Code != 0 {
		err = &status.Error
		return
	}

	neighbours = status.Neighbours
	return
}

type Period int

const (
	Overall Period = 1 + iota
	OneWeek
	OneMonth
	ThreeMonths
	SixMonths
	OneYear
)

var periodStringMap = map[Period]string{
	Overall:     "overall",
	OneWeek:     "7day",
	OneMonth:    "1month",
	ThreeMonths: "3month",
	SixMonths:   "6month",
	OneYear:     "12month"}

func (p Period) String() string {
	return periodStringMap[p]
}

type TopArtists struct {
	User   string `xml:"user,attr"`
	Period Period `xml:"-"`
	Total  int    `xml:"total,attr"`

	Artists []Artist `xml:"artist"`

	// For internal use
	RawPeriod string `xml:"type,attr"`
}

func (top *TopArtists) unmarshalHelper() (err error) {
	for k, v := range periodStringMap {
		if top.RawPeriod == v {
			top.Period = k
			break
		}
	}
	return
}

// Gets a list of the (up to limit) most played artists of a user within a Period.
//
// See http://www.last.fm/api/show/user.getTopArtists.
func (lfm LastFM) GetUserTopArtists(user string, period Period, limit int) (top *TopArtists, err error) {
	bytes, err := lfm.doQuery("user.getTopArtists", map[string]string{
		"user":   user,
		"period": periodStringMap[period],
		"limit":  strconv.Itoa(limit)})
	if err != nil {
		return
	}

	status := lfmStatus{}
	err = xml.Unmarshal(bytes, &status)
	if err != nil {
		return
	}
	if status.Error.Code != 0 {
		err = &status.Error
		return
	}

	top = &status.TopArtists
	err = top.unmarshalHelper()
	return
}

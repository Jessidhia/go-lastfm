package lastfm

import (
	"encoding/xml"
	"strconv"
)

type RecentTracks struct {
	User       string  `xml:"user,attr"`
	Total      int     `xml:"total,attr"`
	Tracks     []Track `xml:"track"`
	NowPlaying *Track  `xml:"-"`
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

	tracks = &status.RecentTracks
	err = tracks.unmarshalHelper()
	return
}

type Tasteometer struct {
	Users   []string `xml:"input>user>name"`
	Score   float32  `xml:"result>score"`
	Artists []string `xml:"result>artists>artist>name"`
}

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

	taste = &status.Tasteometer
	return
}

type Neighbour struct {
	Name  string  `xml:"name"`
	Match float32 `xml:"match"`
}

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
	OneWeek:     "1week",
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

	top = &status.TopArtists
	err = top.unmarshalHelper()
	return
}

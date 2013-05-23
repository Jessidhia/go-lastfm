package lastfm

import (
	"encoding/xml"
	"time"
)

type TrackInfo struct {
	ID             int           `xml:"id"`
	Name           string        `xml:"name"`
	MBID           string        `xml:"mbid"`
	URL            string        `xml:"url"`
	Duration       time.Duration `xml:"-"`
	Listeners      int           `xml:"listeners"`
	TotalPlaycount int           `xml:"playcount"`
	Artist         Artist        `xml:"artist"`

	// Sometimes not present
	Album   *AlbumInfo `xml:"album"`
	TopTags []string   `xml:"toptags>tag>name"`
	Wiki    *Wiki      `xml:"wiki"`

	// Only present if the user parameter isn't empty ("")
	UserPlaycount int  `xml:"userplaycount"`
	UserLoved     bool `xml:"userloved"`

	// For internal use
	RawDuration string `xml:"duration"`
}

func (info *TrackInfo) unmarshalHelper() (err error) {
	info.Duration, err = time.ParseDuration(info.RawDuration + "ms")
	if err != nil {
		return
	}
	if info.Wiki != nil {
		err = info.Wiki.unmarshalHelper()
	}
	return
}

// Gets information for a Track. The user argument can either be empty ("") or specify a last.fm username, in which
// case .UserPlaycount and .UserLoved will be valid in the returned struct. The autocorrect parameter controls whether
// last.fm's autocorrection algorithms should be run on the artist or track names.
//
// The Track struct must specify either the MBID or both Artist.Name and Name.
// Example literals that can be given as the first argument:
//   lastfm.Track{MBID: "mbid"}
//   lastfm.Track{Artist: lastfm.Artist{Name: "Artist"}, Name: "Track"}
//
// See http://www.last.fm/api/show/track.getInfo.
func (lfm LastFM) GetTrackInfo(track Track, user string, autocorrect bool) (info *TrackInfo, err error) {
	query := map[string]string{}
	if autocorrect {
		query["autocorrect"] = "1"
	} else {
		query["autocorrect"] = "0"
	}

	if user != "" {
		query["username"] = user
	}

	if track.MBID != "" {
		query["mbid"] = track.MBID
	} else {
		query["artist"] = track.Artist.Name
		query["track"] = track.Name
	}

	bytes, err := lfm.doQuery("track.getInfo", query)
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

	info = &status.TrackInfo
	err = info.unmarshalHelper()
	return
}

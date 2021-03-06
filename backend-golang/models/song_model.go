package models

type SongInfo struct {
	Title   string `json: "title,omitempty"`
	Artist  string `json: "artist,omitempty"`
	Album   string `json: "album,omitempty"`
	Artwork string `json: "artwork,omitempty"`
	Path    string `json: "path,omitempty"`
}

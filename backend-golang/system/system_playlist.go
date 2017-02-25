package system

import (
	"BrowserMusic/backend-golang/errors"
	"BrowserMusic/backend-golang/models"
	"encoding/base64"
	"encoding/json"
	id3 "github.com/mikkyang/id3-go"
	"io/ioutil"
	"os"
	"strings"
)

type ISystemPlaylist interface {
	Generate() error
	GetAll() (*[]models.PlaylistInfo, error)
	GetSongs(string) (*[]models.SongInfo, error)
	Add(string) error
	AddSong(string, string) error
	Delete(string) error
	DeleteSong(string, string) error
}

type SystemPlaylist struct{}

func NewSystemPlaylist() *SystemPlaylist {
	return &SystemPlaylist{}
}

func (self *SystemPlaylist) Generate() error {
	dir := "../library"

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	plist := models.Playlist{}
	plist.Name = "All Songs"
	plist.Id = 0

	for _, file := range files {
		if file.Name()[len(file.Name())-4:] != ".mp3" {
			continue
		}

		if plist.Artwork == "" {
			path := dir + "/" + file.Name()
			tag, err := id3.Open(path)
			if err == nil {
				artwork := tag.Frame("APIC")
				if artwork != nil {
					plist.Artwork = base64.URLEncoding.EncodeToString(artwork.Bytes())
				}
			}
			tag.Close()
		}
		plist.SongPaths = append(plist.SongPaths, file.Name())
	}

	jsonStr, err := json.Marshal(plist)
	if err != nil {
		return err
	}

	ioutil.WriteFile("../data/All_Songs.playlist", jsonStr, 0644)

	return nil
}

func (self *SystemPlaylist) GetAll() (*[]models.PlaylistInfo, error) {
	dir := "../data"

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var plists []models.PlaylistInfo

	for _, file := range files {
		if file.Name()[len(file.Name())-9:] == ".playlist" {
			plist := models.PlaylistInfo{}

			// get artwork for playlist
			path := dir + "/" + file.Name()
			pfile, err := ioutil.ReadFile(path)
			if err != nil {
				continue // skip file
			}

			var tempPlst models.Playlist
			json.Unmarshal(pfile, &tempPlst)

			plist.Id = tempPlst.Id
			plist.Name = tempPlst.Name
			// plist.Artwork = tempPlst.Artwork

			plists = append(plists, plist)
		}

	}

	return &plists, nil
}

func (self *SystemPlaylist) GetSongs(pStr string) (*[]models.SongInfo, error) {
	file, err := ioutil.ReadFile("../data/" + pStr + ".playlist")
	if err != nil {
		return nil, err
	}

	// get filepaths to songs in playlist
	var playlist models.Playlist
	json.Unmarshal(file, &playlist)

	var songs []models.SongInfo
	for _, path := range playlist.SongPaths {
		// get info from mp3 tags
		tag, err := id3.Open("../library/" + path)
		if err != nil {
			continue
		}

		// title falls back to filename
		title := tag.Title()
		if title == "" {
			title = pStr[:len(pStr)-4]
		}

		var song models.SongInfo
		song.Title = title
		song.Artist = tag.Artist()
		song.Album = tag.Album()
		song.Path = path

		tag.Close()
		songs = append(songs, song)
	}

	return &songs, nil
}

func (self *SystemPlaylist) Add(plistName string) error {
	var playlist models.Playlist
	playlist.Name = plistName

	files, err := ioutil.ReadDir("../data")
	if err != nil {
		return err
	}
	playlist.Id = len(files)

	jsonStr, err := json.Marshal(playlist)
	if err != nil {
		return err
	}

	filename := strings.Replace(plistName, " ", "_", -1)

	ioutil.WriteFile("../data/"+filename+".playlist", jsonStr, 0644)
	return nil
}

func (self *SystemPlaylist) AddSong(pStr string, songPath string) error {
	fullPath := "../data/" + pStr + ".playlist"
	file, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}

	var playlist models.Playlist
	json.Unmarshal(file, &playlist)

	for _, p := range playlist.SongPaths {
		if songPath == p {
			return errors.New("error: file aleady in playlist")
		}
	}
	playlist.SongPaths = append(playlist.SongPaths, songPath)
	jsonStr, err := json.Marshal(playlist)
	if err != nil {
		return err
	}

	ioutil.WriteFile("../data/"+pStr+".playlist", jsonStr, 0644)
	return nil
}

func (self *SystemPlaylist) Delete(plistName string) error {
	fullPath := "../data/" + plistName + ".playlist"
	if err := os.Remove(fullPath); err != nil {
		return err
	}
	return nil
}

func (self *SystemPlaylist) DeleteSong(pStr string, songPath string) error {
	fullPath := "../data/" + pStr + ".playlist"
	file, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}

	var playlist models.Playlist
	json.Unmarshal(file, &playlist)

	for i, p := range playlist.SongPaths {
		if songPath == p {
			if i == len(playlist.SongPaths)-1 {
				playlist.SongPaths = playlist.SongPaths[:i]
			} else {
				playlist.SongPaths = append(playlist.SongPaths[:i], playlist.SongPaths[i+1:]...)
			}
		}
	}

	jsonStr, err := json.Marshal(playlist)
	if err != nil {
		return err
	}

	ioutil.WriteFile("../data/"+pStr+".playlist", jsonStr, 0644)
	return nil
}
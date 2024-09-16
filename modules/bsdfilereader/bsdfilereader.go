package bsdfilereader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"bsvapi/bsddb"
	"bsvapi/types"
)

type BSD_Reader struct {
	Default_path string
	Headers      *[]types.BSD_HeaderGlobal
	Songs        *[]types.BSD_Song
	Sessions     *[]types.BSD_Session
	Session_map  map[string]*types.BSD_Session
	db           *bsddb.DBwrap
	web_ch       *types.WEB_Settings
}

func (reader *BSD_Reader) ReadBSDAsStreams(entry os.DirEntry) ([]string, error) {
	datastreams := []string{}
	file, err := os.Open(reader.Default_path + entry.Name())
	if err != nil {
		return datastreams, err
	}
	defer file.Close()

	f_bytes, err := io.ReadAll(file)
	if err != nil {
		return datastreams, err
	}
	datastreams = strings.Split(string(f_bytes), "\n")

	// cleanup
	for i := 0; i < len(datastreams); i++ {
		datastreams[i] = strings.ReplaceAll(strings.Trim(datastreams[i], "\r"), "\"NaN\"", "-1")
	}
	return datastreams, nil
}

func (reader *BSD_Reader) ParseBSDFromStream(stream string, date string, sid int, bypass_web_notifier bool) error {
	switch true {
	case (strings.Contains(stream, "songID")):
		// song
		song := types.BSD_Song{}
		err := json.Unmarshal([]byte(stream), &song)
		if err != nil {
			return fmt.Errorf("file: %s: JSON error: %v", date, err)
		}
		song.PlayDate = &date
		if reader.Songs == nil {
			return errors.New("could not parse song data. Reader was not initialized")
		}
		err = reader.db.RegisterSongData(&song, sid)
		if err != nil {
			if err.Error() == bsddb.ERR_SONG_DATA_EXISTS {
				// ignore, song exists
			} else {
				return err
			}
		} else {
			log.Printf("Registered new song data: %s - %s (%s; %s)\n", *song.SongArtist, *song.SongName, *song.SongMapper, *song.SongDifficulty)
			if reader.web_ch != nil && reader.web_ch.OnNewSongData != nil && !bypass_web_notifier {
				*reader.web_ch.OnNewSongData <- &song
			}
		}
		*reader.Songs = append(*reader.Songs, song)
	case (strings.Contains(stream, "totalScore")):
		// header
		header := types.BSD_HeaderGlobal{}
		err := json.Unmarshal([]byte(stream), &header)
		if err != nil {
			return err
		}
		header.Date = &date
		if reader.Headers == nil {
			return errors.New("could not parse header data. Reader was not initialized")
		}

		err = reader.db.RegisterBSDHeader(&header, sid)
		if err != nil {
			if err.Error() == bsddb.ERR_HEADER_DATA_EXISTS {
				// ignore, header exists
			} else {
				return err
			}
		} else {
			log.Printf("Registered new header data for session %d\n", sid)
		}

		*reader.Headers = append(*reader.Headers, header)
	default:
		return errors.New("unrecognized stream format for BSD file")
	}
	return nil
}

func (reader *BSD_Reader) FileWatcher() {
	defer time.AfterFunc(time.Second, reader.FileWatcher)
	entries, err := os.ReadDir(reader.Default_path)
	if err != nil {
		log.Println(err)
		return
	}

	for entry_id := 0; entry_id < len(entries); entry_id++ {
		if !strings.HasSuffix(strings.ToLower(entries[entry_id].Name()), ".bsd") {
			// Ignore non .bsd files
			continue
		}

		if strings.HasPrefix(strings.ToLower(entries[entry_id].Name()), "_pbscore") {
			// @TODO: Add _PBScoreGraphs support
			continue
		}

		f_name := entries[entry_id].Name()
		f_full_path := reader.Default_path + f_name
		if reader.Session_map[f_full_path] != nil {
			// session exists, check moddime
			prev_time := reader.Session_map[f_full_path].UpdatedAt
			info, err := entries[entry_id].Info()
			if err != nil {
				log.Println(err)
				continue
			}
			f_moddime := info.ModTime()
			if f_moddime.After(*prev_time) {
				// file was modified
				reader.Session_map[f_full_path].UpdatedAt = &f_moddime

				sid, err := reader.db.UpdateSessionModtime(f_name, reader.Default_path, f_moddime.UnixMilli())
				if err != nil {
					log.Printf("Error updating moddime: %s\n", err.Error())
					continue
				}

				datastreams, err := reader.ReadBSDAsStreams(entries[entry_id])
				if err != nil {
					log.Printf("Error reading streams: %s\n", err.Error())
					continue
				}
				entry_date := strings.Split(entries[entry_id].Name(), ".")[0]

				log.Printf("Processing session: %d\n", sid)
				for stream_id := 0; stream_id < len(datastreams); stream_id++ {
					err = reader.ParseBSDFromStream(datastreams[stream_id], entry_date, sid, false)
					if err != nil {
						if err.Error() == bsddb.ERR_SONG_DATA_EXISTS {
							continue
						}
						log.Printf("Error parsing datastreams: %s\n", err.Error())
					}
				}
			}
		} else {
			// new session detected
			log.Println("New file found: " + f_name)
			datastreams, err := reader.ReadBSDAsStreams(entries[entry_id])
			if err != nil {
				log.Println(err)
				continue
			}
			entry_date := strings.Split(entries[entry_id].Name(), ".")[0]

			info, err := entries[entry_id].Info()
			if err != nil {
				log.Println(err)
				continue
			}
			f_moddime := info.ModTime()
			sid, err := reader.db.CreateNewSession(f_name, reader.Default_path, f_moddime)
			if err != nil {
				if err.Error() == bsddb.ERR_SESSION_EXISTS {
					// ignore, sid is valid
				} else {
					log.Println(err)
					continue
				}
			}

			session := types.BSD_Session{
				Sid:       &sid,
				FileName:  &f_name,
				FilePath:  &reader.Default_path,
				UpdatedAt: &f_moddime,
			}
			*reader.Sessions = append(*reader.Sessions, session)
			reader.Session_map[reader.Default_path+f_name] = &session

			if reader.web_ch != nil && reader.web_ch.OnNewSession != nil {
				*reader.web_ch.OnNewSession <- &session
			}

			log.Printf("Processing session: %d\n", sid)

			for stream_id := 0; stream_id < len(datastreams); stream_id++ {
				err = reader.ParseBSDFromStream(datastreams[stream_id], entry_date, sid, false)
				if err != nil {
					if err.Error() == bsddb.ERR_SONG_DATA_EXISTS {
						continue
					}
					log.Println(err)
				}
			}
		}

	}
}

func (reader *BSD_Reader) Init(db *bsddb.DBwrap, web_channels *types.WEB_Settings) error {
	reader.Headers = &[]types.BSD_HeaderGlobal{}
	reader.Songs = &[]types.BSD_Song{}
	reader.Sessions = &[]types.BSD_Session{}
	reader.Session_map = make(map[string]*types.BSD_Session)
	reader.db = db
	reader.web_ch = web_channels
	switch runtime.GOOS {
	case "windows":
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}

		reader.Default_path = home + "\\AppData\\Roaming\\Beat Savior Data\\"
	default:
		return errors.New("OS is not supported")
	}

	entries, err := os.ReadDir(reader.Default_path)
	if err != nil {
		return err
	}

	for entry_id := 0; entry_id < len(entries); entry_id++ {
		if !strings.HasSuffix(strings.ToLower(entries[entry_id].Name()), ".bsd") {
			// Ignore non .bsd files
			continue
		}

		if strings.HasPrefix(strings.ToLower(entries[entry_id].Name()), "_pbscore") {
			// @TODO: Add _PBScoreGraphs support
			continue
		}

		entry_date := strings.Split(entries[entry_id].Name(), ".")[0]

		info, err := entries[entry_id].Info()
		if err != nil {
			log.Println(err)
			continue
		}
		f_name := entries[entry_id].Name()
		f_moddime := info.ModTime()
		sid, err := reader.db.CreateNewSession(f_name, reader.Default_path, f_moddime)
		if err != nil {
			if err.Error() == bsddb.ERR_SESSION_EXISTS {
				// ignore, sid is valid
			} else {
				log.Println(err)
				continue
			}
		}

		session := types.BSD_Session{
			Sid:       &sid,
			FileName:  &f_name,
			FilePath:  &reader.Default_path,
			UpdatedAt: &f_moddime,
		}
		*reader.Sessions = append(*reader.Sessions, session)
		reader.Session_map[reader.Default_path+f_name] = &session
		log.Printf("Processing session: %d\n", sid)

		datastreams, err := reader.ReadBSDAsStreams(entries[entry_id])
		if err != nil {
			log.Println(err)
			continue
		}

		// inf, _ := entries[entry_id].Info()
		for stream_id := 0; stream_id < len(datastreams); stream_id++ {
			err = reader.ParseBSDFromStream(datastreams[stream_id], entry_date, sid, true)
			if err != nil {
				if err.Error() == bsddb.ERR_SONG_DATA_EXISTS {
					continue
				}
				log.Println(err)
			}
		}
	}
	reader.FileWatcher()
	return nil
}

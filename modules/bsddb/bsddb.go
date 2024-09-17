package bsddb

import (
	T "bsvapi/types"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_FILE_PATH           = "./bsd_storage.sqlite"
	ERR_SESSION_EXISTS     = "failed to create new session. Session exists"
	ERR_SONG_DATA_EXISTS   = "song insert aborted. Song data already exists"
	ERR_HEADER_DATA_EXISTS = "header data insert aborted. Header data already exists"
)

type DBwrap_int interface {
	Init()
	Reset()
	Close()
}

type DBwrap struct {
	DBwrap_int
	db *sql.DB
}

func (db *DBwrap) Init() error {
	if _, err := os.Stat(DB_FILE_PATH); err == nil {
		// Init connect
		var err2 error
		db.db, err2 = sql.Open("sqlite3", DB_FILE_PATH)
		return err2

	} else if errors.Is(err, os.ErrNotExist) {
		// Call DB reset
		return db.Reset()

	} else {
		// Uhhh... Panic I guess
		panic("Database file access failed. Failed to check for file existence. OS error.")

	}
}

func (db *DBwrap) Reset() error {
	// Del old db
	os.Remove(DB_FILE_PATH)

	// recreate file
	var err2 error
	db.db, err2 = sql.Open("sqlite3", DB_FILE_PATH)
	if err2 != nil {
		return err2
	}

	// Reset sequence def.
	resetseq := []string{
		`CREATE TABLE sessions (
			sid        INTEGER PRIMARY KEY AUTOINCREMENT
                       NOT NULL
                       UNIQUE,
			f_name     TEXT    NOT NULL,
			f_path     TEXT    NOT NULL,
			updated_at INTEGER NOT NULL
		)
		STRICT`,
		`CREATE TABLE song_data (
			play_id            INTEGER PRIMARY KEY AUTOINCREMENT
									   UNIQUE
									   NOT NULL,
			sid                INTEGER    REFERENCES sessions (sid) ON DELETE CASCADE
																 ON UPDATE CASCADE
									   NOT NULL,
			songDataType       INT,
			playerID           TEXT,
			songID             TEXT,
			songDifficulty     TEXT,
			songName           TEXT,
			songArtist         TEXT,
			songMapper         TEXT,
			gameMode           TEXT,
			songDifficultyRank INT,
			songSpeed          REAL,
			songStartTime      REAL,
			songDuration       REAL,
			songJumpDistance   REAL,
			trackers           TEXT,
			indexed_at         INTEGER
		)
		STRICT`,
		`CREATE TABLE deep_trackers (
			play_id     INTEGER REFERENCES song_data (play_id) ON DELETE CASCADE
															   ON UPDATE CASCADE
								PRIMARY KEY
								NOT NULL
								UNIQUE,
			noteTracker TEXT
		)
		WITHOUT ROWID,
		STRICT`,
		`CREATE TABLE headers (
			header_id INTEGER PRIMARY KEY AUTOINCREMENT
							  UNIQUE
							  NOT NULL,
			sid       INTEGER     REFERENCES sessions (sid) ON DELETE CASCADE
														ON UPDATE CASCADE,
			data      TEXT
		)
		STRICT`,
		`CREATE INDEX song_data_name_artist_index ON song_data (
			songName ASC,
			songArtist ASC
		)`,
		`CREATE INDEX song_data_indexed_at_index ON song_data (
			indexed_at COLLATE BINARY DESC
		)`,
	}

	// Run reset sequence
	for _, rsql := range resetseq {
		_, err := db.db.Exec(rsql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DBwrap) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *DBwrap) CreateNewSession(f_name string, f_path string, mod_time time.Time) (int, error) {
	sid := -1
	err := db.db.QueryRow(`SELECT sid FROM sessions WHERE f_name=$1 AND f_path=$2`, f_name, f_path).Scan(&sid)
	if err == nil && sid > 0 {
		return sid, errors.New(ERR_SESSION_EXISTS)
	}

	int_time := mod_time.UnixMilli()

	err = db.db.QueryRow(`INSERT INTO sessions(f_name, f_path, updated_at) VALUES ($1, $2, $3) RETURNING sid`, f_name, f_path, int_time).Scan(&sid)
	if err != nil {
		return sid, err
	}

	return sid, nil
}

func (db *DBwrap) RegisterSongData(data *T.BSD_Song, sid int) (int, error) {
	play_id := -1
	if data == nil || data.Trackers == nil {
		return play_id, errors.New("corrupted song data")
	}
	// use tracker data as unique identifier
	play_track_signature_bytes, err := json.Marshal(data.Trackers)
	if err != nil {
		return play_id, errors.New("failed to create track signature")
	}

	pts_str := string(play_track_signature_bytes)

	// perform optimized search
	var existing_plays int
	err = db.db.QueryRow(`SELECT COUNT(play_id) FROM song_data WHERE songName=$1 AND songArtist=$2 AND trackers=$3`,
		data.SongName,
		data.SongArtist,
		pts_str,
	).Scan(&existing_plays)

	if err != nil {
		return play_id, err
	}

	if existing_plays > 0 {
		return play_id, errors.New(ERR_SONG_DATA_EXISTS)
	}

	err = db.db.QueryRow(`INSERT INTO song_data(sid, songDataType, playerID, songID, songDifficulty, songName, songArtist, songMapper, gameMode, songDifficultyRank, songSpeed, songStartTime, songDuration, songJumpDistance, trackers, indexed_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8 ,$9 ,$10, $11, $12, $13, $14, $15, $16) RETURNING play_id`, sid, data.SongDataType, data.PlayerID, data.SongID, data.SongDifficulty, data.SongName, data.SongArtist, data.SongMapper, data.GameMode, data.SongDifficultyRank, data.SongSpeed, data.SongStartTime, data.SongDuration, data.SongJumpDistance, pts_str, time.Now().UnixMilli()).Scan(&play_id)
	if err != nil {
		return play_id, err
	}

	if data.DeepTrackers != nil && data.DeepTrackers.NoteTracker != nil {
		note_tracker_bytes, err := json.Marshal(data.DeepTrackers.NoteTracker)
		if err != nil {
			return play_id, errors.New("failed to create deepTrackers JSON")
		}

		_, err = db.db.Exec(`INSERT INTO deep_trackers(play_id, noteTracker) VALUES ($1, $2)`, play_id, string(note_tracker_bytes))
		if err != nil {
			return play_id, err
		}
	}
	return play_id, nil
}

func (db *DBwrap) RegisterBSDHeader(data *T.BSD_HeaderGlobal, sid int) error {
	if data == nil {
		return errors.New("BSD file header is corrupted")
	}

	header_bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	header_str := string(header_bytes)

	var existing_headers int
	err = db.db.QueryRow(`SELECT COUNT(header_id) FROM headers WHERE sid=$1 AND data=$2`, sid, header_str).Scan(&existing_headers)
	if err != nil {
		return err
	}
	if existing_headers > 0 {
		return errors.New(ERR_HEADER_DATA_EXISTS)
	}
	_, err = db.db.Exec(`INSERT INTO headers(sid, data) VALUES($1, $2)`, sid, header_str)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBwrap) UpdateSessionModtime(f_name string, f_path string, time_ms int64) (int, error) {
	sid := -1
	err := db.db.QueryRow(`SELECT sid FROM sessions WHERE f_name=$1 AND f_path=$2`, f_name, f_path).Scan(&sid)
	if err != nil {
		return sid, err
	}
	_, err = db.db.Exec(`UPDATE sessions SET updated_at=$1 WHERE sid=$2`, time_ms, sid)
	return sid, err
}

func (db *DBwrap) GetLatestSessionData() (*T.BSD_Session, error) {

	sess := T.BSD_Session{}
	sess.Header = &T.BSD_HeaderGlobal{}
	sess.Songs = &[]T.BSD_Song{}

	// get general session data
	err := db.db.QueryRow(`SELECT sid, f_name, f_path, updated_at FROM sessions ORDER BY updated_at DESC LIMIT 1`).Scan(&sess.Sid, &sess.FileName, &sess.FilePath, &sess.UpdatedAtInt)
	if err != nil {
		return nil, err
	}

	// retrieve header data
	var raw_header string
	err = db.db.QueryRow(`SELECT data FROM headers WHERE sid=$1 LIMIT 1`, sess.Sid).Scan(&raw_header)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(raw_header), &sess.Header)
	if err != nil {
		return nil, err
	}

	rows, err := db.db.Query(`SELECT 
	play_id, sid, playerID, 
	songID, songDifficulty, songName,
	songArtist, songMapper, gameMode,
	songDifficultyRank, songSpeed, songDuration,
	songJumpDistance 
	FROM 
	song_data WHERE sid = $1`, sess.Sid)
	if err != nil {
		return &sess, err
	}
	defer rows.Close()

	for rows.Next() {
		song_data := T.BSD_Song{}
		err := rows.Scan(
			&song_data.PlayID, &song_data.SID, &song_data.PlayerID,
			&song_data.SongID, &song_data.SongDifficulty, &song_data.SongName,
			&song_data.SongArtist, &song_data.SongMapper, &song_data.GameMode,
			&song_data.SongDifficultyRank, &song_data.SongSpeed, &song_data.SongDuration,
			&song_data.SongJumpDistance,
		)
		if err != nil {
			return &sess, err
		}
		*sess.Songs = append(*sess.Songs, song_data)
	}

	return &sess, nil

}

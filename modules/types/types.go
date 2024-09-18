package types

import (
	"fmt"
	"math"
	"strings"
	"time"
)

/*
Color representation in RGBA [0;1] format
*/
type RGBA_Color struct {
	R *float64 `json:"r,omitempty"`
	G *float64 `json:"g,omitempty"`
	B *float64 `json:"b,omitempty"`
	A *float64 `json:"a,omitempty"`
}

/*
Beat Savior Global stats and settings tracking object
*/
type BSD_HeaderGlobal struct {
	PlayerID              *string     `json:"playerID,omitempty"`
	AvgCutScore           *int        `json:"averageCutScore,omitempty"`
	BadCutsCount          *int        `json:"badCutsCount,omitempty"`
	ClearedLvlsCount      *int        `json:"clearedLevelsCount,omitempty"`
	FailedLvlsCount       *int        `json:"failedLevelsCount,omitempty"`
	FullComboCount        *int        `json:"fullComboCount,omitempty"`
	GoodCutsCount         *int        `json:"goodCutsCount,omitempty"`
	HandDistanceTravelled *int        `json:"handDistanceTravelled,omitempty"`
	MissedCutsCount       *int        `json:"missedCutsCount,omitempty"`
	PlayedLevelsCount     *int        `json:"playedLevelsCount,omitempty"`
	TotalScore            *int        `json:"totalScore,omitempty"`
	TimePlayed            *float64    `json:"timePlayed,omitempty"`
	SaberAColor           *RGBA_Color `json:"saberAColor,omitempty"`
	SaberBColor           *RGBA_Color `json:"saberBColor,omitempty"`
	LightAColor           *RGBA_Color `json:"lightAColor,omitempty"`
	LightBColor           *RGBA_Color `json:"lightBColor,omitempty"`
	ObstacleColor         *RGBA_Color `json:"obstacleColor,omitempty"`

	// Custom Field
	Date *string `json:"date,omitempty"`
}

/*
Beat Savior Song Definition
*/
type BSD_Song struct {
	SongDataType       *int              `json:"songDataType,omitempty"`
	PlayerID           *string           `json:"playerID,omitempty"`
	SongID             *string           `json:"songID,omitempty"`
	SongDifficulty     *string           `json:"songDifficulty,omitempty"`
	SongName           *string           `json:"songName,omitempty"`
	SongArtist         *string           `json:"songArtist,omitempty"`
	SongMapper         *string           `json:"songMapper,omitempty"`
	GameMode           *string           `json:"gameMode,omitempty"`
	SongDifficultyRank *int              `json:"songDifficultyRank,omitempty"`
	SongSpeed          *float64          `json:"songSpeed,omitempty"`
	SongStartTime      *float64          `json:"songStartTime,omitempty"`
	SongDuration       *float64          `json:"songDuration,omitempty"`
	SongJumpDistance   *float64          `json:"songJumpDistance,omitempty"`
	Trackers           *BSD_Trackers     `json:"trackers,omitempty"`
	DeepTrackers       *BSD_DeepTrackers `json:"deepTrackers,omitempty"`

	// Custom Field
	IndexedAtInt *int    `json:"indexed_at_int,omitempty"`
	SearchQuery  *string `json:"search_query,omitempty"`
	PlayDate     *string `json:"play_date,omitempty"`
	PlayID       *int    `json:"play_id,omitempty"`
	SID          *int    `json:"sid,omitempty"`
}

type BSD_Trackers struct {
	HitTracker        *BSD_HitTracker        `json:"hitTracker,omitempty"`
	AccuracyTracker   *BSD_AccuracyTracker   `json:"accuracyTracker,omitempty"`
	ScoreTracker      *BSD_ScoreTracker      `json:"scoreTracker,omitempty"`
	WinTracker        *BSD_WinTracker        `json:"winTracker,omitempty"`
	DistanceTracker   *BSD_DistanceTracker   `json:"distanceTracker,omitempty"`
	ScoreGraphTracker *BSD_ScoreGraphTracker `json:"scoreGraphTracker,omitempty"`
}

type BSD_HitTracker struct {
	LeftNoteHit      *int `json:"leftNoteHit,omitempty"`
	RightNoteHit     *int `json:"rightNoteHit,omitempty"`
	BombHit          *int `json:"bombHit,omitempty"`
	MaxCombo         *int `json:"maxCombo,omitempty"`
	NumberOfWallsHit *int `json:"nbOfWallHit,omitempty"`
	Misses           *int `json:"miss,omitempty"`
	MissedNotes      *int `json:"missedNotes,omitempty"`
	BadCuts          *int `json:"badCuts,omitempty"`
	LeftMiss         *int `json:"leftMiss,omitempty"`
	LeftBadCuts      *int `json:"leftBadCuts,omitempty"`
	RightMiss        *int `json:"rightMiss,omitempty"`
	RightBadCuts     *int `json:"rightBadCuts,omitempty"`
}

type BSD_AccuracyTracker struct {
	AccRight              *float64   `json:"accRight,omitempty"`
	AccLeft               *float64   `json:"accLeft,omitempty"`
	AverageAcc            *float64   `json:"averageAcc,omitempty"`
	LeftSpeed             *float64   `json:"leftSpeed,omitempty"`
	RightSpeed            *float64   `json:"rightSpeed,omitempty"`
	AverageSpeed          *float64   `json:"averageSpeed,omitempty"`
	LeftHighestSpeed      *float64   `json:"leftHighestSpeed,omitempty"`
	RightHighestSpeed     *float64   `json:"rightHighestSpeed,omitempty"`
	LeftPreswing          *float64   `json:"leftPreswing,omitempty"`
	RightPreswing         *float64   `json:"rightPreswing,omitempty"`
	AveragePreswing       *float64   `json:"averagePreswing,omitempty"`
	LeftPostswing         *float64   `json:"leftPostswing,omitempty"`
	RightPostswing        *float64   `json:"rightPostswing,omitempty"`
	AveragePostswing      *float64   `json:"averagePostswing,omitempty"`
	LeftTimeDependence    *float64   `json:"leftTimeDependence,omitempty"`
	RightTimeDependence   *float64   `json:"rightTimeDependence,omitempty"`
	AverageTimeDependence *float64   `json:"averageTimeDependence,omitempty"`
	LeftAverageCut        *[]float64 `json:"leftAverageCut,omitempty"`
	RightAverageCut       *[]float64 `json:"rightAverageCut,omitempty"`
	AverageCut            *[]float64 `json:"averageCut,omitempty"`
	GridAcc               *[]float64 `json:"gridAcc,omitempty"`
	GridCut               *[]float64 `json:"gridCut,omitempty"`
}

type BSD_ScoreTracker struct {
	RawScore                  *int      `json:"rawScore,omitempty"`
	Score                     *int      `json:"score,omitempty"`
	PersonalBest              *int      `json:"personalBest,omitempty"`
	RawRatio                  *float64  `json:"rawRatio,omitempty"`
	ModifiedRatio             *float64  `json:"modifiedRatio,omitempty"`
	PersonalBestRawRatio      *float64  `json:"personalBestRawRatio,omitempty"`
	PersonalBestModifiedRatio *float64  `json:"personalBestModifiedRatio,omitempty"`
	ModifiersMultiplier       *float64  `json:"modifiersMultiplier,omitempty"`
	Modifiers                 *[]string `json:"modifiers,omitempty"`
}

type BSD_WinTracker struct {
	Won            *bool    `json:"won,omitempty"`
	Rank           *string  `json:"rank,omitempty"`
	EndTime        *float64 `json:"endTime,omitempty"`
	NumberOfPauses *int     `json:"nbOfPause,omitempty"`
}

type BSD_DistanceTracker struct {
	RightSaber *float64 `json:"rightSaber,omitempty"`
	LeftSaber  *float64 `json:"leftSaber,omitempty"`
	RightHand  *float64 `json:"rightHand,omitempty"`
	LeftHand   *float64 `json:"leftHand,omitempty"`
}

type BSD_ScoreGraphTracker struct {
	Graph *map[string]float64 `json:"graph,omitempty"`
}

type BSD_DeepTrackers struct {
	NoteTracker *BSD_NoteTracker `json:"noteTracker,omitempty"`
}

type BSD_NoteTracker struct {
	Notes *[]BSD_Note `json:"notes,omitempty"`
}

type BSD_Note struct {
	NoteType         *int       `json:"noteType,omitempty"`
	NoteDirection    *int       `json:"noteDirection,omitempty"`
	Index            *int       `json:"index,omitempty"`
	ID               *int       `json:"id,omitempty"`
	Time             *float64   `json:"time,omitempty"`
	CutType          *int       `json:"cutType,omitempty"`
	Multiplier       *int       `json:"multiplier,omitempty"`
	Score            *[]int     `json:"score,omitempty"`
	NoteCenter       *[]float64 `json:"noteCenter,omitempty"`
	NoteRotation     *[]float64 `json:"noteRotation,omitempty"`
	TimeDeviation    *float64   `json:"timeDeviation,omitempty"`
	Speed            *float64   `json:"speed,omitempty"`
	Preswing         *float64   `json:"preswing,omitempty"`
	Postswing        *float64   `json:"postswing,omitempty"`
	DistanceToCenter *float64   `json:"distanceToCenter,omitempty"`
	CutPoint         *[]float64 `json:"cutPoint,omitempty"`
	SaberDir         *[]float64 `json:"saberDir,omitempty"`
	CutNormal        *[]float64 `json:"cutNormal,omitempty"`
	TimeDependence   *float64   `json:"timeDependence,omitempty"`
}

type BSD_Session struct {
	Sid       *int       `json:"sid,omitempty"`
	FileName  *string    `json:"f_name,omitempty"`
	FilePath  *string    `json:"f_path,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// API
	UpdatedAtInt *int              `json:"updated_at_int,omitempty"`
	Header       *BSD_HeaderGlobal `json:"header,omitempty"`
	Songs        *[]BSD_Song       `json:"songs,omitempty"`
}

type WEB_Settings struct {
	OnNewSession  *chan *BSD_Session
	OnNewSongData *chan *BSD_Song
	OnReady       *chan struct{}
}

/*
Returns a color representation in hex format,
i.e. #RRGGBBAA
*/
func (color *RGBA_Color) ToHex() string {
	out := strings.Builder{}
	out.WriteString("#")
	for _, clr_val := range []*float64{color.R, color.G, color.B, color.A} {
		if clr_val != nil {
			// Convert the [0;1] range to [0;255], clamp it, and return as int
			x := int(min(255, max(0, math.Floor(*clr_val*255.0))))
			out.WriteString(fmt.Sprintf("%X", x))
		} else {
			out.WriteString("00")
		}
	}
	return out.String()
}

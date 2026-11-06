package models

import (
	"time"

	"github.com/photoview/photoview/api/guessgame"
)

// GuessGame is one game of the 4D-guesser, played on one device of an account.
//
// Fork-specific: not part of upstream Photoview.
type GuessGame struct {
	Model
	// PublicID is the random id the API uses, so games cannot be reached by counting.
	PublicID string `gorm:"size:32;not null;uniqueIndex"`
	UserID   int    `gorm:"not null;index"`
	User     User   `gorm:"constraint:OnDelete:CASCADE;"`
	// DeviceKey is a random id of the browser that started the game, used to avoid repeats.
	DeviceKey       string     `gorm:"size:64;not null;index"`
	RoundCount      int        `gorm:"not null"`
	Difficulty      string     `gorm:"size:16;not null"`
	DistanceScaleKm float64    `gorm:"not null"`
	TimeScaleMonths float64    `gorm:"not null"`
	MinRating       int        `gorm:"not null"`
	TotalScore      int        `gorm:"not null;default:0"`
	FinishedAt      *time.Time `gorm:"index"`
	// Abandoned is set when the device starts another game, or when the game has been idle too long.
	Abandoned bool `gorm:"not null;default:false;index"`
}

func (GuessGame) TableName() string {
	return "guess_games"
}

// Scales returns the scoring scales the game was started with.
func (g *GuessGame) Scales() guessgame.Scales {
	return guessgame.Scales{DistanceKm: g.DistanceScaleKm, TimeMonths: g.TimeScaleMonths}
}

// MaxScore is the highest total the game can reach.
func (g *GuessGame) MaxScore() int {
	return g.RoundCount * guessgame.MaxRoundScore
}

// GuessRound is one photo of a GuessGame. The answer is copied from the photo's EXIF data
// when the game starts, so later changes to the photo do not change the score.
//
// Fork-specific: not part of upstream Photoview.
type GuessRound struct {
	Model
	PublicID   string     `gorm:"size:32;not null;uniqueIndex"`
	GameID     int        `gorm:"not null;uniqueIndex:idx_guess_rounds_game_round"`
	Game       *GuessGame `gorm:"constraint:OnDelete:CASCADE;"`
	RoundIndex int        `gorm:"not null;uniqueIndex:idx_guess_rounds_game_round"`
	MediaID    *int       `gorm:"index"`
	Media      *Media     `gorm:"constraint:OnDelete:SET NULL;"`

	AnswerLatitude  float64 `gorm:"not null"`
	AnswerLongitude float64 `gorm:"not null"`
	AnswerYear      int     `gorm:"not null"`
	AnswerMonth     int     `gorm:"not null"`

	GuessLatitude  *float64
	GuessLongitude *float64
	GuessYear      *int
	GuessMonth     *int
	DistanceKm     *float64
	MonthsOff      *int
	LocationScore  *int
	TimeScore      *int
	GuessedAt      *time.Time
}

func (GuessRound) TableName() string {
	return "guess_rounds"
}

// Guessed reports whether the round has a stored guess.
func (r *GuessRound) Guessed() bool {
	return r.GuessedAt != nil
}

// GuessRoundResult is the guess of a round, the answer and the score.
//
// Fork-specific: not part of upstream Photoview.
type GuessRoundResult struct {
	Guess         *Coordinates
	GuessYear     int
	GuessMonth    int
	Answer        *Coordinates
	AnswerYear    int
	AnswerMonth   int
	DistanceKm    float64
	MonthsOff     int
	LocationScore int
	TimeScore     int
	Score         int
}

// Result returns the stored result of the round, or nil if it has not been guessed yet.
func (r *GuessRound) Result() *GuessRoundResult {
	if !r.Guessed() || r.GuessLatitude == nil || r.GuessLongitude == nil || r.GuessYear == nil ||
		r.GuessMonth == nil || r.DistanceKm == nil || r.MonthsOff == nil || r.LocationScore == nil || r.TimeScore == nil {
		return nil
	}

	return &GuessRoundResult{
		Guess:         &Coordinates{Latitude: *r.GuessLatitude, Longitude: *r.GuessLongitude},
		GuessYear:     *r.GuessYear,
		GuessMonth:    *r.GuessMonth,
		Answer:        &Coordinates{Latitude: r.AnswerLatitude, Longitude: r.AnswerLongitude},
		AnswerYear:    r.AnswerYear,
		AnswerMonth:   r.AnswerMonth,
		DistanceKm:    *r.DistanceKm,
		MonthsOff:     *r.MonthsOff,
		LocationScore: *r.LocationScore,
		TimeScore:     *r.TimeScore,
		Score:         *r.LocationScore + *r.TimeScore,
	}
}

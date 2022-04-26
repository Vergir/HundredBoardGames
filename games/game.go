package games

import (
	"fmt"
	"math"
)

const (
	RATING_ID_WILSON = 0
)

type Game struct {
	GeekId              uint
	PrimaryTitle        string
	Titles              []string
	Year                int16
	Description         string
	PictureUrl          string
	MinPlayers          uint8
	MaxPlayers          uint8
	CommunityNumPlayers []numPlayersPoll
	MinPlaytime         uint
	MaxPlaytime         uint
	MinAge              uint8
	CommunityMinAge     []minAgePoll
	LanguageDependence  []langDepPoll
	Tags                []tag
	AvgRating           float64
	BayesRating         float64
	RatingNumVotes      uint
	Ratings             []float64
	AvgWeight           float64
	WeightNumVotes      uint
	Counters            gameCounters
}

type numPlayersPoll struct {
	NumPlayers          uint8
	VotedBest           uint
	VotedRecommended    uint
	VotedNotRecommended uint
}

type minAgePoll struct {
	MinAge   uint8
	NumVotes uint
}

type langDepPoll struct {
	Level    uint8
	NumVotes uint
}

type tag struct {
	Type string
	Id   uint
}

type gameCounters struct {
	Owned    uint
	Trading  uint
	Wanting  uint
	Wishing  uint
	Comments uint
}

func (game Game) String() string {
	return fmt.Sprintf("[%10v] %50v | Avg: %4.3v | Bayes: %4.3v | NumVotes: (%v)\n",
		game.GeekId, game.PrimaryTitle, game.AvgRating, game.BayesRating, game.RatingNumVotes)
}

func wilson(p float64, n uint) float64 {
	const wilsonZ = 1.96
	const wilsonZSquared = 3.8416

	flN := float64(n)

	a := p + wilsonZSquared/(2*flN)
	pre_b := p*(1-p) + wilsonZSquared/(4*flN)
	b := wilsonZ * math.Sqrt(pre_b/flN)
	c := 1 + wilsonZSquared/flN

	result := (a - b) / c

	return result
}

func calcWilsonRating(votersRating float64, votersCount uint) float64 {
	wilsonP := votersRating / 100
	wilsonN := votersCount

	result := wilson(wilsonP, wilsonN) * 100

	return result
}

func (game *Game) UpdateAlgoRatings() {
	if cap(game.Ratings) < 256 {
		game.Ratings = make([]float64, 256)
	}
	wilson_rating := calcWilsonRating(game.AvgRating, game.RatingNumVotes)
	game.Ratings[RATING_ID_WILSON] = wilson_rating
}

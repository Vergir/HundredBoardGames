package games

import (
	"fmt"
	"math"
)

type Game struct {
	Title        string
	Year         int
	GeekRating   float64
	VotersRating float64
	VotersCount  uint
	AlgoRating   float64
}

func NewGame(title string, year int, geekRating float64, votersRating float64, votersCount uint) Game {
	algoRating := calcAlgoRating(votersRating, votersCount)

	game := Game{title, year, geekRating, votersRating, votersCount, algoRating}

	return game
}

func PrintGame(game Game) {
	fmt.Printf("%50v (%4v) | Geek: %4.3v | Algo: %4.3v | Vote: %4.3v (%v)\n",
		game.Title, game.Year, game.GeekRating, game.AlgoRating, game.VotersRating, game.VotersCount)
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

func calcAlgoRating(votersRating float64, votersCount uint) float64 {
	wilsonP := votersRating / 100
	wilsonN := votersCount

	result := wilson(wilsonP, wilsonN) * 100

	return result
}

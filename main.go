package main

import (
	"fmt"
	"hundred-board-games/games"
)

func main() {
	gamesList, _ := games.ReadGamesFromStorage()

	topGeek := games.GetTopGames(gamesList, games.SORT_BY_GEEK_RATING, 50)
	for _, e := range topGeek {
		games.PrintGame(e)
	}

	fmt.Println()
	fmt.Println()
	fmt.Println()

	topAlgo := games.GetTopGames(gamesList, games.SORT_BY_ALGO_RATING, 50)
	for _, e := range topAlgo {
		games.PrintGame(e)
	}
}

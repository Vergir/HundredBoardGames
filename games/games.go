package games

import "sort"

const (
	SORT_BY_GEEK_RATING = "geek"
	SORT_BY_ALGO_RATING = "algo"
)

func formSorter(games []Game, sortType string) func(a, b int) bool {
	var sorter func(a, b int) bool

	switch sortType {
	case SORT_BY_GEEK_RATING:
		sorter = func(a, b int) bool {
			return games[a].GeekRating > games[b].GeekRating
		}
	case SORT_BY_ALGO_RATING:
		sorter = func(a, b int) bool {
			return games[a].AlgoRating > games[b].AlgoRating
		}
	}

	return sorter
}

func GetTopGames(games []Game, sortType string, topN uint) []Game {
	sort.Slice(games, formSorter(games, sortType))

	return games[:topN]
}

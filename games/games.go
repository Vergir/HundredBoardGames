package games

import "sort"

const (
	SORT_BY_GEEK_RATING = iota
	SORT_BY_ALGO_RATING = iota
)

func formSorter(games []Game, sortType int) func(a, b int) bool {
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

func GetTopGames(games []Game, sortType int, topN uint) []Game {
	sort.Slice(games, formSorter(games, sortType))

	return games[:topN]
}

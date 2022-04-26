package games

import "sort"

func formSorter(games []Game, ratingId uint8) func(a, b int) bool {
	var sorter func(a, b int) bool

	switch ratingId {
	case RATING_ID_WILSON:
		sorter = func(a, b int) bool {
			return games[a].Ratings[RATING_ID_WILSON] > games[b].Ratings[RATING_ID_WILSON]
		}
	}

	return sorter
}

func GetTopGames(games []Game, ratingId uint8, topN uint) []Game {
	sort.Slice(games, formSorter(games, ratingId))

	return games[:topN]
}

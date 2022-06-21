package pages

import (
	"hundred-board-games/code/server/paths"
	"hundred-board-games/code/utils"
)

var INDEX_PAGE = newPage(
	"index", "index", "Hundred Board Games", paths.PAGE_INDEX,
	[]string{utils.StaticJs("common.js")},
	[]string{"index.css"},
)

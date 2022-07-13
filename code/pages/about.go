package pages

import (
	"hundred-board-games/code/server/paths"
)

var ABOUT_PAGE = newPage(
	"about", "about", "about", paths.PAGE_ABOUT,
	[]string{},
	[]string{"about.css"},
)

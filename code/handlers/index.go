package handlers

import (
	"fmt"
	"hundred-board-games/code/pages"
	"hundred-board-games/code/templates"
	"net/http"
)

type indexPageTemplateProps struct {
	Intro        string
	About        string
	TopPrompt    string
	FilterPrompt string
	ButtonTop    string
	ButtonFilter string
}

func prepareIndexPageProps() indexPageTemplateProps {
	intro := "A website that ranks board games. The algorithm looks at how people rate games and applies math to do the ranking."
	about := "More info here"
	buttonTop := "Top games"
	buttonFilter := "Seach tool"
	props := indexPageTemplateProps{
		Intro:        intro,
		About:        about,
		TopPrompt:    "See the list of highest ranked games",
		FilterPrompt: "Or search for a particular game",
		ButtonTop:    buttonTop,
		ButtonFilter: buttonFilter,
	}

	return props
}

func HandleIndexPage(request *http.Request, headers http.Header) (string, error) {
	indexPageProps := prepareIndexPageProps()

	response, err := templates.RenderPage(pages.INDEX_PAGE, indexPageProps)
	if err != nil {
		fmt.Println(err)
	}

	return response, nil
}

package handlers

import (
	"fmt"
	"html/template"
	"hundred-board-games/code/i18n"
	"hundred-board-games/code/pages"
	"hundred-board-games/code/templates"
	"net/http"
)

type aboutPageTemplateProps struct {
	TechnicalDetailsHtml template.HTML
}

func prepareAboutPageProps() aboutPageTemplateProps {
	tokensUrls := map[string]string{
		"blogpost": "https://www.evanmiller.org/how-not-to-sort-by-average-rating.html",
		"bgg":      "https://boardgamegeek.com/",
		"github":   "https://github.com/",
	}
	tokensExpansions := make(map[string]string)
	for token, url := range tokensUrls {
		tokensExpansions[token] = fmt.Sprintf("<a href=\"%v\" rel=\"noopener noreferrer\" target=\"_blank\">%v</a>", url, i18n.Get(pages.ABOUT_PAGE.LangSection, token))
	}
	technicalDetails := i18n.Expand(i18n.Get(pages.ABOUT_PAGE.LangSection, "paragraph_technical_details"), tokensExpansions)

	props := aboutPageTemplateProps{
		TechnicalDetailsHtml: template.HTML(technicalDetails),
	}

	return props
}

func HandleAboutPage(request *http.Request, headers http.Header) (string, error) {
	props := prepareAboutPageProps()
	response, err := templates.RenderPage(pages.ABOUT_PAGE, props)
	if err != nil {
		fmt.Println(err)
	}

	return response, nil
}

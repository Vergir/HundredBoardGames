package about

import (
	"fmt"
	"html/template"
	"hundred-board-games/code/i18n"
)

type aboutPageTemplateProps struct {
	TechnicalDetailsHtml template.HTML
}

func buildPageProps() aboutPageTemplateProps {
	tokensUrls := map[string]string{
		"blogpost": "https://www.evanmiller.org/how-not-to-sort-by-average-rating.html",
		"bgg":      "https://boardgamegeek.com/",
		"github":   "https://github.com/",
	}
	tokensExpansions := make(map[string]string)
	for token, url := range tokensUrls {
		tokensExpansions[token] = fmt.Sprintf("<a href=\"%v\" rel=\"noopener noreferrer\" target=\"_blank\">%v</a>", url, i18n.Get(ENDPOINT.I18nSection, token))
	}
	technicalDetails := i18n.Expand(i18n.Get(ENDPOINT.I18nSection, "paragraph_technical_details"), tokensExpansions)

	props := aboutPageTemplateProps{
		TechnicalDetailsHtml: template.HTML(technicalDetails),
	}

	return props
}

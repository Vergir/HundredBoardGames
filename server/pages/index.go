package pages

import (
	"html/template"
	"hundred-board-games/server"
	"strings"
)

type chapter struct {
	Title string
	Url   string
}

type indexPageData struct {
	PageTitle string
	Chapters  []chapter
}

func RenderIndexPage() (string, error) {
	tmpl, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		return "", err
	}

	data := indexPageData{
		PageTitle: "TOP GAMES LIST",
		Chapters: []chapter{
			{Title: "TOP 100", Url: makeChapterUrl(server.PATH_LIST)},
		},
	}

	var stringBuilder strings.Builder

	err = tmpl.Execute(&stringBuilder, data)
	if err != nil {
		return "", err
	}

	return stringBuilder.String(), nil
}

func makeChapterUrl(path string) string {
	return path
}

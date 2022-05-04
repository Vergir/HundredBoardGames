package server

import (
	"errors"
	"fmt"
	"html/template"
	"hundred-board-games/server/pages"
	"net/http"
	"strings"
	"time"
)

type props struct {
	Global globalProps
	Page   pageProps
}

type globalProps struct {
	PageTitle string
	Now       uint
}

type pageProps any

var templates *template.Template = template.Must(template.ParseGlob("templates/*"))

func GetAndRenderTemplate(page pages.Page, pageProps any) (string, error) {
	template := templates.Lookup(page.TemplateName + ".tmpl")
	if template == nil {
		return "", errors.New("no template found")
	}

	templateProps := props{
		Global: globalProps{
			PageTitle: page.Title,
			Now:       uint(time.Now().Unix()),
		},
		Page: pageProps,
	}

	var stringBuilder strings.Builder

	err := template.Execute(&stringBuilder, templateProps)
	if err != nil {
		return "", err
	}

	return stringBuilder.String(), nil
}

func AddStatic() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func AddHandler(url string, handler func(request *http.Request) string) {
	path := "/" + url + "/"
	if url == pages.INDEX_PAGE.Url {
		path = "/"
	}

	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		//debug
		templates = template.Must(template.ParseGlob("templates/*"))

		fmt.Fprint(writer, handler(request))
	})
}

func Start() error {
	return http.ListenAndServe(":80", nil)
}

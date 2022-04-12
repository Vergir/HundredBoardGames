package server

import (
	"errors"
	"fmt"
	"html/template"
	"hundred-board-games/server/pages"
	"net/http"
	"strings"
)

var templates *template.Template = template.Must(template.ParseGlob("templates/*"))

func GetAndRenderTemplate(templateName string, templateData any) (string, error) {
	template := templates.Lookup(templateName + ".tmpl")
	if template == nil {
		return "", errors.New("no template found")
	}

	var stringBuilder strings.Builder

	err := template.Execute(&stringBuilder, templateData)
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
		fmt.Fprint(writer, handler(request))
	})
}

func Start() error {
	return http.ListenAndServe(":80", nil)
}

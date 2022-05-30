package templates

import (
	"errors"
	"fmt"
	"html/template"
	"hundred-board-games/code/pages"
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
	JsPaths   []string
}

type pageProps any

var templates *template.Template = template.Must(template.ParseGlob("templates/*"))

func RenderPage(page pages.Page, pageProps any) (string, error) {
	template := templates.Lookup(page.TemplateName + ".tmpl")
	if template == nil {
		return "", errors.New("no template found")
	}

	jsPaths := make([]string, len(page.JsPaths))
	for i, jsPath := range page.JsPaths {
		if strings.Contains(jsPath, "static") {
			jsPath += fmt.Sprint("?v=", time.Now().Unix())
		}
		jsPaths[i] = jsPath
	}

	templateProps := props{
		Global: globalProps{
			PageTitle: page.Title,
			Now:       uint(time.Now().Unix()),
			JsPaths:   jsPaths,
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

func RenderCustom(templateName string, props any) (string, error) {
	template := templates.Lookup(templateName + ".tmpl")
	if template == nil {
		return "", errors.New("no template found")
	}

	var stringBuilder strings.Builder

	err := template.Execute(&stringBuilder, props)
	if err != nil {
		return "", err
	}

	return stringBuilder.String(), nil
}

func Reload() {
	templates = template.Must(template.ParseGlob("templates/*"))
}

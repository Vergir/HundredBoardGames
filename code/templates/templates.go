package templates

import (
	"errors"
	"fmt"
	"html/template"
	"hundred-board-games/code/i18n"
	"hundred-board-games/code/website"
	"strings"
	"time"
)

type props struct {
	Global globalProps
	Page   pageProps
}

type globalProps struct {
	Lang        string
	SharedI18n  map[string]string
	PageI18n    map[string]string
	CurrentYear uint
	JsPaths     []string
	CssPaths    []string
}

type pageProps any

var templates *template.Template = template.Must(template.ParseGlob("templates/*"))

func RenderEndpoint(endpoint *website.Endpoint, pageProps any) (string, error) {
	template := templates.Lookup(endpoint.TemplateName)
	if template == nil {
		return "", errors.New("no template found")
	}

	jsPaths := make([]string, len(endpoint.JsPaths))
	for i, jsPath := range endpoint.JsPaths {
		if strings.Contains(jsPath, "static") {
			jsPath += fmt.Sprint("?v=", time.Now().Unix())
		}
		jsPaths[i] = jsPath
	}

	cssPaths := make([]string, len(endpoint.CssPaths))
	//debug
	for i, cssPath := range endpoint.CssPaths {
		cssPath += fmt.Sprint("?v=", time.Now().Unix())
		cssPaths[i] = cssPath
	}

	templateProps := props{
		Global: globalProps{
			Lang:        string(i18n.GetCurrentLocale()),
			SharedI18n:  i18n.GetSection("shared"),
			PageI18n:    i18n.GetSection(endpoint.I18nSection),
			CurrentYear: uint(time.Now().Year()),
			JsPaths:     jsPaths,
			CssPaths:    cssPaths,
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

func Reload() {
	templates = template.Must(template.ParseGlob("templates/*"))
}

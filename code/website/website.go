package website

import (
	"net/http"
	"strings"
)

type RequestHandler func(request *http.Request, headers http.Header) (string, error)

type Endpoint struct {
	CodeName     string
	TemplateName string
	I18nSection  string
	Url          string
	JsPaths      []string
	CssPaths     []string
}

func NewEndpoint(codeName string, staticPaths ...string) *Endpoint {
	jsPaths, cssPaths := formStaticsSlices(staticPaths)

	endpoint := Endpoint{
		CodeName:     codeName,
		TemplateName: codeName + ".tmpl",
		I18nSection:  codeName,
		Url:          codeName,
		JsPaths:      jsPaths,
		CssPaths:     cssPaths,
	}

	return &endpoint
}

func NewComplexEndpoint(codeName, templateName, i18nSection, url string, jsPaths []string, cssPaths []string) *Endpoint {
	endpoint := Endpoint{
		CodeName:     codeName,
		TemplateName: templateName,
		I18nSection:  i18nSection,
		Url:          url,
		JsPaths:      jsPaths,
		CssPaths:     cssPaths,
	}

	return &endpoint
}

func formStaticsSlices(staticPaths []string) (jsPaths []string, cssPaths []string) {
	jsPaths = make([]string, 0)
	cssPaths = make([]string, 0)
	for _, filePath := range staticPaths {
		dotIndex := strings.LastIndexByte(filePath, '.')
		if dotIndex == -1 {
			continue
		}
		switch filePath[dotIndex+1:] {
		case "js":
			jsPaths = append(jsPaths, filePath)
		case "css":
			cssPaths = append(cssPaths, filePath)
		}
	}

	return jsPaths, cssPaths
}

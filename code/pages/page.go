package pages

type Page struct {
	CodeName     string
	TemplateName string
	LangSection  string
	Url          string
	JsPaths      []string
	CssPaths     []string
}

func newPage(codeName, templateName, langSection, url string, jsPaths []string, cssPaths []string) Page {
	return Page{codeName, templateName, langSection, url, jsPaths, cssPaths}
}

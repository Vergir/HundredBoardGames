package pages

type Page struct {
	CodeName     string
	TemplateName string
	Title        string
	Url          string
	JsPaths      []string
	CssPaths     []string
}

func newPage(codeName, templateName, title, url string, jsPaths []string, cssPaths []string) Page {
	return Page{codeName, templateName, title, url, jsPaths, cssPaths}
}

var ALL_PAGES = []Page{
	TOP_PAGE,
	INDEX_PAGE,
}

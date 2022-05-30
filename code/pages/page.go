package pages

type Page struct {
	CodeName     string
	TemplateName string
	Title        string
	Url          string
	JsPaths      []string
}

func newPage(codeName, templateName, title, url string, jsPaths ...string) Page {
	return Page{codeName, templateName, title, url, jsPaths}
}

var ALL_PAGES = []Page{
	TOP_PAGE,
	INDEX_PAGE,
}

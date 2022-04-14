package pages

type Page struct {
	CodeName     string
	TemplateName string
	Title        string
	Url          string
}

type PageProps interface {
	GetFinalTemplateProps() any
	SetPageTitle(pageTitle string)
}

func newPage(codeName, templateName, title, url string) Page {
	return Page{codeName, templateName, title, url}
}

var ALL_PAGES = []Page{
	TOP_PAGE,
	INDEX_PAGE,
}
package pages

type Page struct {
	Name         string
	Url          string
	TemplateName string
}

func newPage(name, url, templateName string) Page {
	return Page{name, url, templateName}
}

var ALL_PAGES = []Page{
	LIST_PAGE,
	INDEX_PAGE,
}

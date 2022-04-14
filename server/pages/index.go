package pages

type chapter struct {
	Title string
	Url   string
}

type indexPageTemplateProps struct {
	PageTitle string
	Chapters  []chapter
}

var INDEX_PAGE = newPage("index", "index", "Hundred Board Games", "/")

func (props *indexPageTemplateProps) SetPageTitle(title string) {
	props.PageTitle = title
}

func (props *indexPageTemplateProps) GetFinalTemplateProps() any {
	return *props
}

func PrepareIndexPageProps() PageProps {
	chapters := make([]chapter, 0)
	for _, page := range ALL_PAGES {
		chapters = append(chapters, chapter{Title: page.Title, Url: page.Url})
	}

	props := indexPageTemplateProps{
		Chapters: chapters,
	}

	return &props
}

package pages

type chapter struct {
	Title string
	Url   string
}

type indexPageTemplateProps struct {
	Chapters []chapter
}

var INDEX_PAGE = newPage("index", "index", "Hundred Board Games", "/")

func PrepareIndexPageProps() indexPageTemplateProps {
	chapters := make([]chapter, 0)
	for _, page := range ALL_PAGES {
		chapters = append(chapters, chapter{Title: page.Title, Url: page.Url})
	}

	props := indexPageTemplateProps{
		Chapters: chapters,
	}

	return props
}

package pages

type chapter struct {
	Title string
	Url   string
}

type indexPageData struct {
	PageTitle string
	Chapters  []chapter
}

var INDEX_PAGE = newPage("index", "/", "index")

func PrepareIndexPageData() any {
	//TODO: proper chapters
	data := indexPageData{
		PageTitle: "TOP GAMES LIST",
		Chapters: []chapter{
			{Title: "TOP 100", Url: ""},
		},
	}

	return data
}

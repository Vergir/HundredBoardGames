package server

import (
	"fmt"
	"hundred-board-games/code/handlers"
	"hundred-board-games/code/pages"
	"hundred-board-games/code/templates"
	"net/http"
)

func AddStatic() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func AddHandler(url string, handler handlers.Handler) {
	path := "/" + url + "/"
	if url == pages.INDEX_PAGE.Url {
		path = "/"
	}

	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		//debug
		templates.Reload()

		response, err := handler(request, writer.Header())
		if err != nil {
			fmt.Fprint(writer, "Sorry")
			fmt.Println(err)
		}

		fmt.Fprint(writer, response)
	})
}

func Start() error {
	return http.ListenAndServe(":80", nil)
}

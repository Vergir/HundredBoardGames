package server

import (
	"fmt"
	"hundred-board-games/code/i18n"
	"hundred-board-games/code/templates"
	"hundred-board-games/code/website"
	"hundred-board-games/code/website/index"
	"net/http"
)

func AddStatic() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func AddEndpointHandler(endpoint *website.Endpoint, requestHandler website.RequestHandler) {
	path := "/" + endpoint.Url + "/"
	if endpoint.CodeName == index.ENDPOINT.CodeName {
		path = "/"
	}

	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		i18n.Init()
		//debug
		reloadResources()

		response, err := requestHandler(request, writer.Header())
		if err != nil {
			fmt.Fprint(writer, "Sorry")
			fmt.Println(err)
		}

		fmt.Fprint(writer, response)
	})
}

//use only for debug
func reloadResources() {
	i18n.LoadLocale(i18n.LOCALE_EN_GB)
	templates.Reload()
}

func Start() error {
	return http.ListenAndServe(":80", nil)
}

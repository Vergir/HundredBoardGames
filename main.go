package main

import (
	"hundred-board-games/code/server"
	"hundred-board-games/code/website/about"
	"hundred-board-games/code/website/index"
	"hundred-board-games/code/website/top"
	"hundred-board-games/code/website/top/gamesextras"
)

func main() {
	server.AddEndpointHandler(index.ENDPOINT, index.HANDLER)
	server.AddEndpointHandler(about.ENDPOINT, about.HANDLER)
	server.AddEndpointHandler(top.ENDPOINT, top.HANDLER)
	server.AddEndpointHandler(gamesextras.ENDPOINT, gamesextras.HANDLER)

	server.AddStatic()

	server.Start()
}

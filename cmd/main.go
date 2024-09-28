package main

import (
	"MusicLibraryAPI/internal/app"

	_ "github.com/lib/pq"
)

func main() {

	app.StartServer()
}

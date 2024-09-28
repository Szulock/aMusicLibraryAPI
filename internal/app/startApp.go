package app

import (
	"MusicLibraryAPI/internal/transport"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var logger = logrus.New()

func StartServer() {
	http.HandleFunc("/songs", transport.GetSongs)

	http.HandleFunc("/songCr", transport.CreateSong)
	http.HandleFunc("/songUp/", transport.UpdateSong)
	http.HandleFunc("/songDel/", transport.DeleteSong)
	http.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Infof("Запуск сервера на порту %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal("Не получилось запустить сервер:", err)
	}
}

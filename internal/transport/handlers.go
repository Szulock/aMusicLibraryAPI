package transport

import (
	"MusicLibraryAPI/internal/database"
	"MusicLibraryAPI/internal/models"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var db = database.GetDB()

func GetSongs(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "БД не создана", http.StatusInternalServerError)
		return
	}

	var songs []models.Song
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}
	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum
	query := `SELECT * FROM songs LIMIT $1 OFFSET $2`
	err = db.Select(&songs, query, limitNum, offset)
	if err != nil {
		logger.Error("Не удалось получить список песен:", err)
		http.Error(w, "Не удалось получить список песен", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// GetSongText godoc
// @Summary Получить текст песни по ID
// @Description Получает текст песни по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {array} string
// @Failure 404 {object} models.ErrorResponse
// @Router /songText/{id} [get]
func CreateSong(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "БД не создана", http.StatusInternalServerError)
		return
	}

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Некорректный ответ", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO songs (group_name, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.QueryRow(query, song.GroupName, song.Song, song.ReleaseDate, song.Text, song.Link).Scan(&song.ID)
	if err != nil {
		logger.Error("Не удалось создать песню:", err)
		http.Error(w, "Не удалось создать песню", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

func UpdateSong(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "БД не создана", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Path[len("/songUp/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный id песни", http.StatusBadRequest)
		return
	}

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Некорректный ответ", http.StatusBadRequest)
		return
	}

	query := `UPDATE songs SET group_name = $1, song = $2, release_date = $3, text = $4, link = $5 WHERE id = $6`
	_, err = db.Exec(query, song.GroupName, song.Song, song.ReleaseDate, song.Text, song.Link, id)
	if err != nil {
		logger.Error("Не удалось обновить песню:", err)
		http.Error(w, "Не удалось обновить песню", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteSong(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "БД не создана", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Path[len("/songDel/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный номер песни", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM songs WHERE id = $1`
	_, err = db.Exec(query, id)
	if err != nil {
		logger.Error("Не удалось удалить песню:", err)
		http.Error(w, "Не удалось удалить песню", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

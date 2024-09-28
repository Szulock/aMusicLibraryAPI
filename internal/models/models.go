package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group_name"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

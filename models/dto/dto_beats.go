package dto

import (
	"time"

	"github.com/google/uuid"
)

type BeatAndUser struct {
	ID 			uuid.UUID 	`json:"id"`
	Username 	string    	`json:"username"`
	FileURL     string 		`json:"file_url"`
}

type FullBeatAndUser struct {
	ID 			uuid.UUID 	`json:"id"`
	Username    string	`json:"username"`
	Title       string   	`json:"title"`
	Description string		`json:"description"`
	Genre       string 		`json:"genre"`
	Tags        string		`json:"tags"`
	FileURL     string 		`json:"file_url"`
	FileSize    int64 		`json:"file_size"`
	CreatedAt   time.Time 	`json:"created_at"`
}

type BeatAndLike struct {
	ID uuid.UUID `json:"id"`
	Title string `json:"title"`
	Username string `json:"username"`
	FileURL string `json:"file_url"`
	Likes uint64 `json:"likes"`
}
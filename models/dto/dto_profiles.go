package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProfileInformation struct {
	ID		 string `json:"id"`
	Username string `json:"username"`
	Region   string `json:"region"`
	Discord  string `json:"discord"`
	Tracks   int64  `json:"tracks"`
	Likes    int64  `json:"likes"`
}

type BeatByUser struct {
	ID          uuid.UUID  `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	Tags        string 	  `json:"tags"`
	FileURL     string    `json:"file_url"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
	Likes       int64     `json:"likes"`
}
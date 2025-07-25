package dto

import "github.com/google/uuid"

type BeatAndUser struct {
	ID uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FileURL     string 		`json:"file_url"`
}
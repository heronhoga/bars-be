package responses

import (
	"time"

	"github.com/google/uuid"
)

type GetBeatResponses struct {
	ID          uuid.UUID 	`json:"id"`
	UserID      uuid.UUID 	`json:"user_id"`
	Title       string   	`json:"title"`
	Description string		`json:"description"`
	Genre       string 		`json:"genre"`
	Tags        string		`json:"tags"`
	FileURL     string 		`json:"file_url"`
	FileSize    int64 		`json:"file_size"`
	CreatedAt   time.Time 	`json:"created_at"`
}
package requests

type CreateBeatRequest struct {
	Title       string `json:"title" validate:"required,max=20"`
	Description string `json:"description" validate:"required,max=200"`
	Genre       string `json:"genre" validate:"required"`
	Tags        string `json:"tags" validate:"required"`
}

type DeleteBeatRequest struct {
	BeatID string `json:"beat_id" validate:"required"`
}

type EditBeatRequest struct {
	Title       string `json:"title" validate:"required,max=20"`
	Description string `json:"description" validate:"required,max=200"`
	Genre       string `json:"genre" validate:"required"`
	Tags        string `json:"tags" validate:"required"`
}
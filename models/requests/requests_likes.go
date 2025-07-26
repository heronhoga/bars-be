package requests

type LikeRequest struct {
	BeatID string `json:"beat_id" validate:"required"`
}
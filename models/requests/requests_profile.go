package requests

type EditProfileRequest struct {
	ID       string `json:"id" validate:"required"`
	Region   string `json:"region" validate:"required"`
	Discord  string `json:"discord" validate:"required"`
}
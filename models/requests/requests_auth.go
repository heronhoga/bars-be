package requests

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Region   string `json:"region" validate:"required"`
	Discord  string `json:"discord" validate:"required"`
}
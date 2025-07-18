package entities

type UserEntity struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Region   string `json:"region"`
}
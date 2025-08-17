package dto

type ProfileInformation struct {
	Username string `json:"username"`
	Region string `json:"region"`
	Discord string `json:"discord"`
	Tracks int64 `json:"tracks"`
	Likes int64 `json:"likes"`
}
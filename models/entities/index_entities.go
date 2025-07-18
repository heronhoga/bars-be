package entities

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Region   string `gorm:"type:varchar(100)"`
	Discord  string `gorm:"type:varchar(255)"`

	Beats      []Beat     `gorm:"foreignKey:UserID"`
	LikedBeats []LikedBeat `gorm:"foreignKey:UserID"`
}

type Beat struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null"`
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Title       string `gorm:"not null"`
	Description string
	Genre       string `gorm:"type:varchar(100)"`
	Tags        string
	FileURL     string `gorm:"not null"`
	FileSize    int64
	CreatedAt   time.Time

	Likes []LikedBeat `gorm:"foreignKey:BeatID"`
}

type LikedBeat struct {
	ID     uint `gorm:"primaryKey"`
	BeatID uint `gorm:"not null"`
	Beat   Beat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserID uint `gorm:"not null"`
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
}

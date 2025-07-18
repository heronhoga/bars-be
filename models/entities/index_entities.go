package entities

import "time"

type UserEntity struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Region   string `gorm:"type:varchar(100)"`
	Discord  string `gorm:"type:varchar(255)"`

	Beats      []BeatEntity      `gorm:"foreignKey:UserID"`
	LikedBeats []LikedBeatEntity `gorm:"foreignKey:UserID"`
}

type BeatEntity struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null"`
	User        UserEntity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Title       string `gorm:"not null"`
	Description string
	Genre       string `gorm:"type:varchar(100)"`
	Tags        string
	FileURL     string `gorm:"not null"`
	FileSize    int64
	CreatedAt   time.Time

	Likes []LikedBeatEntity `gorm:"foreignKey:BeatID"`
}

type LikedBeatEntity struct {
	ID     uint `gorm:"primaryKey"`
	BeatID uint `gorm:"not null"`
	Beat   BeatEntity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserID uint `gorm:"not null"`
	User   UserEntity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
}

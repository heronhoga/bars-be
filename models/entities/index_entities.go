package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Region   string    `gorm:"type:varchar(100)"`
	Discord  string    `gorm:"type:varchar(255)"`

	Beats      []Beat     `gorm:"foreignKey:UserID"`
	LikedBeats []LikedBeat `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}


type Beat struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Title       string    `gorm:"not null"`
	Description string
	Genre       string `gorm:"type:varchar(100)"`
	Tags        string
	FileURL     string    `gorm:"not null"`
	FileSize    int64
	CreatedAt   time.Time

	Likes []LikedBeat `gorm:"foreignKey:BeatID"`
}

func (b *Beat) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}


type LikedBeat struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	BeatID uuid.UUID `gorm:"type:uuid;not null"`
	Beat   Beat      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
}

func (lb *LikedBeat) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	lb.ID = id
	return nil
}


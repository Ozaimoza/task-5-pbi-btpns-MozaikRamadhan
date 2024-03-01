package models

import (
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	// ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	// CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	// UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	// Relation to Photo
	Photos []PhotoModel `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type PhotoModel struct {
	gorm.Model
	// ID        uint      `json:"id" gorm:"primaryKey"`
	Title    string `json:"title" gorm:"not null"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photoUrl" gorm:"not null"`
	UserID   uint   `json:"userId" gorm:"not null"`
	// CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	// UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

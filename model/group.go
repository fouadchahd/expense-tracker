package model

import "time"

type Group struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	Code         string        `json:"code" gorm="unique"`
	CreatedAt    time.Time     `json:"created_at"`
	Users        []User        `json:"-" gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Transactions []Transaction `json:"-" gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

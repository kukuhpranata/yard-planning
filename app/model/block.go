package model

import (
	"time"
)

type Block struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	YardID int    `gorm:"not null" json:"yard_id"`
	Name   string `gorm:"type:varchar(50);not null;uniqueIndex:idx_yard_block" json:"name"`

	Slots int `gorm:"not null" json:"slots"` // Panjang
	Rows  int `gorm:"not null" json:"rows"`  // Lebar
	Tiers int `gorm:"not null" json:"tiers"` // Tinggi

	Yard Yard `gorm:"foreignKey:YardID;references:ID" json:"yard,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

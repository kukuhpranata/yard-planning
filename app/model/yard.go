package model

import (
	"time"
)

type Yard struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"type:varchar(100);unique;not null" json:"name"`
	Location string `gorm:"type:varchar(255)" json:"location"`

	Blocks []Block `gorm:"foreignKey:YardID" json:"blocks,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

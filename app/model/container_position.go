package model

import (
	"time"
)

type ContainerPosition struct {
	ID              int    `gorm:"primaryKey" json:"id"`
	ContainerNumber string `gorm:"type:varchar(20);unique;not null" json:"container_number"`

	BlockID    int `gorm:"not null" json:"block_id"`
	SlotNumber int `gorm:"not null;uniqueIndex:idx_position" json:"slot_number"`
	RowNumber  int `gorm:"not null;uniqueIndex:idx_position" json:"row_number"`
	TierNumber int `gorm:"not null;uniqueIndex:idx_position" json:"tier_number"`

	ContainerSize   string `gorm:"type:varchar(5);not null" json:"container_size"`
	ContainerHeight string `gorm:"type:varchar(5);not null" json:"container_height"`
	ContainerType   string `gorm:"type:varchar(50);not null" json:"container_type"`

	ContainerStatus string    `gorm:"type:varchar(20);not null" json:"container_status"`
	ArrivalDate     time.Time `gorm:"type:timestamp with time zone" json:"arrival_date"`

	YardPlanID *int `gorm:"null" json:"yard_plan_id,omitempty"`

	Block Block `gorm:"foreignKey:BlockID;references:ID" json:"block,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

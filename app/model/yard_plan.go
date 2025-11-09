package model

import (
	"time"
)

type YardPlan struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	BlockID  int    `gorm:"not null" json:"block_id"`
	PlanName string `gorm:"type:varchar(255);not null" json:"plan_name"`

	SlotStart int `gorm:"not null" json:"slot_start"`
	SlotEnd   int `gorm:"not null" json:"slot_end"`
	RowStart  int `gorm:"not null" json:"row_start"`
	RowEnd    int `gorm:"not null" json:"row_end"`

	ContainerSize   string `gorm:"type:varchar(5);not null" json:"container_size"`   // '20ft', '40ft'
	ContainerHeight string `gorm:"type:varchar(5);not null" json:"container_height"` // '8.6ft', '9.6ft'
	ContainerType   string `gorm:"type:varchar(50);not null" json:"container_type"`

	PriorityStackingDirection string `gorm:"type:varchar(50)" json:"priority_stacking_direction"`
	IsActive                  bool   `gorm:"not null;default:true" json:"is_active"`

	Block Block `gorm:"foreignKey:BlockID;references:ID" json:"block,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

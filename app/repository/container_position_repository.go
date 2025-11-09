package repository

import (
	"errors"
	"yard-planning/app/model"

	"gorm.io/gorm"
)

type ContainerPositionRepository interface {
	Save(db *gorm.DB, position *model.ContainerPosition) error
	FindByContainerNumber(db *gorm.DB, positionResult *model.ContainerPosition, containerNumber string) error
	Delete(db *gorm.DB, containerID int) error

	CheckPositionAvailability(db *gorm.DB, blockID, row, tier int, slotNumbers []int) (int64, error)
	IsStackedAbove(db *gorm.DB, blockID, slot, row, tier int) (bool, error)
}

type ContainerPositionRepositoryImpl struct {
}

func NewContainerPositionRepository() ContainerPositionRepository {
	return &ContainerPositionRepositoryImpl{}
}

func (r *ContainerPositionRepositoryImpl) Save(db *gorm.DB, position *model.ContainerPosition) error {
	query := `INSERT INTO container_positions (
		container_number, block_id, slot_number, row_number, tier_number, 
		container_size, container_height, container_type, container_status, 
		arrival_date, yard_plan_id, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result := db.Exec(query,
		position.ContainerNumber, position.BlockID, position.SlotNumber, position.RowNumber, position.TierNumber,
		position.ContainerSize, position.ContainerHeight, position.ContainerType, position.ContainerStatus,
		position.ArrivalDate, position.YardPlanID, position.CreatedAt, position.UpdatedAt,
	)

	if result.RowsAffected == 0 {
		return errors.New("failed to insert container position")
	}
	return result.Error
}

func (r *ContainerPositionRepositoryImpl) FindByContainerNumber(db *gorm.DB, positionResult *model.ContainerPosition, containerNumber string) error {
	err := db.Raw("SELECT * FROM container_positions WHERE container_number = ?", containerNumber).Scan(positionResult).Error

	if errors.Is(err, gorm.ErrRecordNotFound) || positionResult.ID == 0 {
		return errors.New("container not found at any position")
	}
	return err
}

func (r *ContainerPositionRepositoryImpl) Delete(db *gorm.DB, containerID int) error {
	result := db.Exec("DELETE FROM container_positions WHERE id = ?", containerID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("container position not found or already deleted")
	}
	return nil
}

func (r *ContainerPositionRepositoryImpl) CheckPositionAvailability(db *gorm.DB, blockID, row, tier int, slotNumbers []int) (int64, error) {
	var count int64

	query := `
		SELECT COUNT(id) FROM container_positions
		WHERE block_id = ? AND row_number = ? AND tier_number = ? AND slot_number IN (?)`

	result := db.Raw(query, blockID, row, tier, slotNumbers).Scan(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (r *ContainerPositionRepositoryImpl) IsStackedAbove(db *gorm.DB, blockID, slot, row, tier int) (bool, error) {
	var count int64

	query := `
		SELECT COUNT(id) FROM container_positions
		WHERE block_id = ? 
		  AND slot_number = ? 
		  AND row_number = ? 
		  AND tier_number > ?`

	result := db.Raw(query, blockID, slot, row, tier).Scan(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

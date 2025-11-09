package repository

import (
	"errors"
	"yard-planning/app/model"

	"gorm.io/gorm"
)

type YardPlanRepository interface {
	Save(db *gorm.DB, plan *model.YardPlan) error
	FindByID(db *gorm.DB, planResult *model.YardPlan, planID int) error
	Delete(db *gorm.DB, planID int) error

	FindActivePlansByBlock(db *gorm.DB, plans *[]model.YardPlan, blockID int) error

	FindOverlappingPlans(db *gorm.DB, plans *[]model.YardPlan, newPlan *model.YardPlan) error
	FindApplicablePlan(db *gorm.DB, blockID, slot, row int, size, height, cType string) (*model.YardPlan, error)
}

type YardPlanRepositoryImpl struct {
}

func NewYardPlanRepository() YardPlanRepository {
	return &YardPlanRepositoryImpl{}
}

func (r *YardPlanRepositoryImpl) Save(db *gorm.DB, plan *model.YardPlan) error {
	query := `INSERT INTO yard_plans (
		block_id, plan_name, slot_start, slot_end, row_start, row_end, 
		container_size, container_height, container_type, priority_stacking_direction, 
		is_active, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result := db.Exec(query,
		plan.BlockID, plan.PlanName, plan.SlotStart, plan.SlotEnd, plan.RowStart, plan.RowEnd,
		plan.ContainerSize, plan.ContainerHeight, plan.ContainerType, plan.PriorityStackingDirection,
		plan.IsActive, plan.CreatedAt, plan.UpdatedAt,
	)

	if result.RowsAffected == 0 {
		return errors.New("failed to insert yard plan")
	}
	return result.Error
}

func (r *YardPlanRepositoryImpl) FindByID(db *gorm.DB, planResult *model.YardPlan, planID int) error {
	err := db.Raw("SELECT * FROM yard_plans WHERE id = ?", planID).Scan(planResult).Error

	if errors.Is(err, gorm.ErrRecordNotFound) || planResult.ID == 0 {
		return errors.New("yard plan not found")
	}
	return err
}

func (r *YardPlanRepositoryImpl) Delete(db *gorm.DB, planID int) error {
	result := db.Exec("DELETE FROM yard_plans WHERE id = ?", planID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("yard plan not found or already deleted")
	}
	return nil
}

func (r *YardPlanRepositoryImpl) FindActivePlansByBlock(db *gorm.DB, plans *[]model.YardPlan, blockID int) error {
	query := `
		SELECT * FROM yard_plans 
		WHERE block_id = ? AND is_active = TRUE
		ORDER BY id ASC`

	err := db.Raw(query, blockID).Scan(plans).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (r *YardPlanRepositoryImpl) FindOverlappingPlans(db *gorm.DB, plans *[]model.YardPlan, newPlan *model.YardPlan) error {
	query := `
		SELECT * FROM yard_plans AS old_plan
		WHERE 
			old_plan.block_id = ? AND
			(old_plan.id <> ?) AND -- Hindari membandingkan diri sendiri saat UPDATE
			
			-- Slot Overlap Check
			old_plan.slot_end >= ? AND 
			old_plan.slot_start <= ? AND
			
			-- Row Overlap Check
			old_plan.row_end >= ? AND 
			old_plan.row_start <= ? AND
			
			-- Conflict Check: Hanya jika spesifikasi kontainer BERBEDA.
			(old_plan.container_size <> ? OR
			 old_plan.container_height <> ? OR
			 old_plan.container_type <> ?)
	`

	err := db.Raw(query,
		newPlan.BlockID, newPlan.ID,
		newPlan.SlotStart, newPlan.SlotEnd,
		newPlan.RowStart, newPlan.RowEnd,
		newPlan.ContainerSize, newPlan.ContainerHeight, newPlan.ContainerType,
	).Scan(plans).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (r *YardPlanRepositoryImpl) FindApplicablePlan(db *gorm.DB, blockID, slot, row int, size, height, cType string) (*model.YardPlan, error) {
	var planResult model.YardPlan

	query := `
		SELECT * FROM yard_plans 
		WHERE 
			block_id = ? AND 
			is_active = TRUE AND 
			
			-- Cek Koordinat: Slot harus berada di antara start dan end
			slot_start <= ? AND 
			slot_end >= ? AND 
			row_start <= ? AND 
			row_end >= ? AND 
			
			-- Cek Spesifikasi Kontainer
			container_size = ? AND 
			container_height = ? AND 
			container_type = ?
		LIMIT 1
	`

	err := db.Raw(query,
		blockID,
		slot, slot,
		row, row,
		size, height, cType,
	).Scan(&planResult).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || planResult.ID == 0 {
			return nil, nil
		}
		return nil, err
	}

	return &planResult, nil
}

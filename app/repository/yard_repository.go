package repository

import (
	"errors"
	"yard-planning/app/model"

	"gorm.io/gorm"
)

type YardRepository interface {
	SaveYard(db *gorm.DB, yard *model.Yard) error
	FindYardByName(db *gorm.DB, yardResult *model.Yard, name string) error
	FindYardByID(db *gorm.DB, yardResult *model.Yard, yardID int) error
	DeleteYard(db *gorm.DB, yardID int) error

	FindBlockByNameAndYardID(db *gorm.DB, blockResult *model.Block, blockName string, yardID int) error
	FindBlocksByYardID(db *gorm.DB, blockResults *[]model.Block, yardID int) error
	FindBlockByID(db *gorm.DB, blockResult *model.Block, blockID int) error
}

type YardRepositoryImpl struct {
}

func NewYardRepository() YardRepository {
	return &YardRepositoryImpl{}
}

func (r *YardRepositoryImpl) SaveYard(db *gorm.DB, yard *model.Yard) error {
	query := `INSERT INTO yards (name, location, created_at, updated_at) VALUES (?, ?, ?, ?)`
	result := db.Exec(query, yard.Name, yard.Location, yard.CreatedAt, yard.UpdatedAt)

	if result.RowsAffected == 0 {
		return errors.New("failed to insert yard")
	}
	return result.Error
}

func (r *YardRepositoryImpl) FindYardByName(db *gorm.DB, yardResult *model.Yard, name string) error {
	// GORM's .Scan() can map raw query result to struct.
	err := db.Raw("SELECT * FROM yards WHERE name = ?", name).Scan(yardResult).Error

	if errors.Is(err, gorm.ErrRecordNotFound) || yardResult.ID == 0 {
		// GORM Raw doesn't always return ErrRecordNotFound, check if struct is empty
		return errors.New("yard not found")
	}
	return err
}

func (r *YardRepositoryImpl) FindYardByID(db *gorm.DB, yardResult *model.Yard, yardID int) error {
	err := db.Raw("SELECT * FROM yards WHERE id = ?", yardID).Scan(yardResult).Error

	if errors.Is(err, gorm.ErrRecordNotFound) || yardResult.ID == 0 {
		return errors.New("yard not found")
	}
	return err
}

func (r *YardRepositoryImpl) DeleteYard(db *gorm.DB, yardID int) error {
	result := db.Exec("DELETE FROM yards WHERE id = ?", yardID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("yard not found or already deleted")
	}
	return nil
}

func (r *YardRepositoryImpl) FindBlockByNameAndYardID(db *gorm.DB, blockResult *model.Block, blockName string, yardID int) error {
	err := db.Raw("SELECT * FROM blocks WHERE name = ? AND yard_id = ?", blockName, yardID).Scan(blockResult).Error

	if errors.Is(err, gorm.ErrRecordNotFound) || blockResult.ID == 0 {
		return errors.New("block not found in yard")
	}
	return err
}

func (r *YardRepositoryImpl) FindBlocksByYardID(db *gorm.DB, blockResults *[]model.Block, yardID int) error {
	err := db.Raw("SELECT * FROM blocks WHERE yard_id = ?", yardID).Scan(blockResults).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (r *YardRepositoryImpl) FindBlockByID(db *gorm.DB, blockResult *model.Block, blockID int) error {
	query := `SELECT * FROM blocks WHERE id = ?`
	err := db.Raw(query, blockID).Scan(blockResult).Error

	if errors.Is(err, gorm.ErrRecordNotFound) || blockResult.ID == 0 {
		return errors.New("block not found")
	}

	if err != nil {
		return err
	}

	return nil
}

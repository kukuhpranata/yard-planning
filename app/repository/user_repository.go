package repository

import (
	"errors"
	"yard-planning/app/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Save(db *gorm.DB, user *model.User) error
	FindByEmail(db *gorm.DB, userResult *model.User, email string) error
	FindById(db *gorm.DB, userResult *model.User, userId int) error
	Update(db *gorm.DB, user *model.User) error
	Delete(db *gorm.DB, userId int) error
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r UserRepositoryImpl) Save(db *gorm.DB, user *model.User) error {
	query := `INSERT INTO users (name, email, password, created_at, updated_at) VALUES (?,?,?,?,?,?)`
	result := db.Exec(query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if result.RowsAffected == 0 {
		return errors.New("failed to insert")
	}
	return nil
}

func (r UserRepositoryImpl) FindByEmail(db *gorm.DB, userResult *model.User, email string) error {
	err := db.Raw("SELECT * from users where email = ?", email).Scan(&userResult).Error
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepositoryImpl) FindById(db *gorm.DB, userResult *model.User, userId int) error {
	err := db.Raw("SELECT * from users where id = ?", userId).Scan(&userResult).Error
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepositoryImpl) Update(db *gorm.DB, user *model.User) error {
	result := db.Exec("UPDATE users SET name = ?, password = ?, updated_at = ? where id = ?", user.Name, user.Password, user.UpdatedAt, user.Id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("data not found")
	}
	return nil
}

func (r UserRepositoryImpl) Delete(db *gorm.DB, userId int) error {
	result := db.Exec("DELETE FROM users where id = ?", userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return errors.New("data not found")
	}
	return nil
}

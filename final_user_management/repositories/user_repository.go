package repositories

import (
	"github.com/umwaribenie/final_user_management/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(params models.GetAllUsersRequest) ([]models.User, int64, error)
	FindByID(id string) (*models.User, error)
	FindBySlug(slug string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(id string, user *models.User) error
	Delete(id string) error
	UpdatePassword(id string, password string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll(params models.GetAllUsersRequest) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if params.From != "" && params.To != "" {
		query = query.Where("created_at BETWEEN ? AND ?", params.From, params.To)
	}
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR username LIKE ?", search, search, search, search)
	}
	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.PageSize == 0 {
		params.PageSize = 10
	}
	if params.PageNumber == 0 {
		params.PageNumber = 1
	}
	offset := (params.PageNumber - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindBySlug(slug string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(id string, user *models.User) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) UpdatePassword(id string, password string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password", password).Error
}

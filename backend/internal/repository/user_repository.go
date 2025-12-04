package repository

import (
	"context"

	"github.com/easyhire/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	ListUsers(ctx context.Context, filter UserFilter) ([]models.User, int64, error)
}

type UserFilter struct {
	Role      models.UserRole
	Email     string
	Search    string
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	return &user, err
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) ListUsers(ctx context.Context, filter UserFilter) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.User{})
	
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}
	if filter.Search != "" {
		query = query.Where("email ILIKE ? OR name ILIKE ?", 
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	
	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}
	
	// Apply sorting
	if filter.SortBy != "" {
		order := filter.SortBy
		if filter.SortOrder != "" {
			order = order + " " + filter.SortOrder
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}
	
	err = query.Find(&users).Error
	return users, total, err
}

package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/Mluciano1401/go-api/internal/application"
	"github.com/Mluciano1401/go-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:150;uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time
}

func (userModel) TableName() string { return "users" }

type UserModelForMigration = userModel

func toDomain(m userModel) *domain.User {
	return &domain.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
	}
}

func fromDomain(u *domain.User) userModel {
	return userModel{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ application.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	m := fromDomain(user)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var m userModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m userModel
	err := r.db.WithContext(ctx).First(&m, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	var models []userModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	users := make([]domain.User, 0, len(models))
	for _, m := range models {
		users = append(users, *toDomain(m))
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	m := fromDomain(user)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&userModel{}, "id = ?", id).Error
}

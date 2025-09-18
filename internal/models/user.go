package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password  string         `json:"-" gorm:"not null" validate:"required,min=6"`
	FirstName string         `json:"first_name" gorm:"not null" validate:"required"`
	LastName  string         `json:"last_name" gorm:"not null" validate:"required"`
	Phone     string         `json:"phone" validate:"required"`
	Role      string         `json:"role" gorm:"default:user" validate:"oneof=user owner admin"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Courts   []Court   `json:"courts,omitempty" gorm:"foreignKey:OwnerID"`
	Bookings []Booking `json:"bookings,omitempty" gorm:"foreignKey:UserID"`
	Reviews  []Review  `json:"reviews,omitempty" gorm:"foreignKey:UserID"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

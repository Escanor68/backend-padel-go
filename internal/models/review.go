package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CourtID   uint           `json:"court_id" gorm:"not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Rating    int            `json:"rating" gorm:"not null" validate:"required,min=1,max=5"`
	Comment   string         `json:"comment" gorm:"type:text" validate:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Court Court `json:"court,omitempty" gorm:"foreignKey:CourtID"`
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"required"`
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment *string `json:"comment,omitempty"`
}

type ReviewResponse struct {
	ID        uint      `json:"id"`
	CourtID   uint      `json:"court_id"`
	UserID    uint      `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Información del usuario
	User UserInfo `json:"user"`
}

type UserInfo struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

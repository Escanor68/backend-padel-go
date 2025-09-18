package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CourtID    uint           `json:"court_id" gorm:"not null"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	Date       time.Time      `json:"date" gorm:"type:date;not null"`
	StartTime  string         `json:"start_time" gorm:"not null" validate:"required"`
	EndTime    string         `json:"end_time" gorm:"not null" validate:"required"`
	Status     string         `json:"status" gorm:"default:pending" validate:"oneof=pending confirmed cancelled completed"`
	TotalPrice float64        `json:"total_price" gorm:"type:decimal(10,2)" validate:"required,min=0"`
	Notes      string         `json:"notes"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Court Court `json:"court,omitempty" gorm:"foreignKey:CourtID"`
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type CreateBookingRequest struct {
	CourtID   uint   `json:"court_id" validate:"required"`
	Date      string `json:"date" validate:"required"` // formato: "2024-03-20"
	StartTime string `json:"start_time" validate:"required"` // formato: "10:00"
	EndTime   string `json:"end_time" validate:"required"`   // formato: "11:00"
	Notes     string `json:"notes"`
}

type BookingResponse struct {
	ID         uint      `json:"id"`
	CourtID    uint      `json:"court_id"`
	UserID     uint      `json:"user_id"`
	Date       time.Time `json:"date"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	Status     string    `json:"status"`
	TotalPrice float64   `json:"total_price"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// Información adicional
	Court CourtInfo `json:"court"`
	User  UserInfo  `json:"user"`
}

type CourtInfo struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	PricePerHour float64 `json:"price_per_hour"`
	Surface     string  `json:"surface"`
	HasLighting bool    `json:"has_lighting"`
	IsIndoor    bool    `json:"is_indoor"`
}

type UserInfo struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type GetBookingsRequest struct {
	UserID  *uint   `json:"user_id,omitempty"`
	CourtID *uint   `json:"court_id,omitempty"`
	Date    *string `json:"date,omitempty"` // formato: "2024-03-20"
	Status  *string `json:"status,omitempty"`
}

type TimeSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Available bool   `json:"available"`
	Price     float64 `json:"price,omitempty"`
}

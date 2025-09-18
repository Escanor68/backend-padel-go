package models

import (
	"time"

	"gorm.io/gorm"
)

type Court struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"not null" validate:"required"`
	Address         string         `json:"address" gorm:"not null" validate:"required"`
	Latitude        float64        `json:"latitude" gorm:"type:decimal(10,7)" validate:"required"`
	Longitude       float64        `json:"longitude" gorm:"type:decimal(10,7)" validate:"required"`
	PricePerHour    float64        `json:"price_per_hour" gorm:"type:decimal(10,2)" validate:"required,min=0"`
	Description     string         `json:"description"`
	ImageURL        string         `json:"image_url"`
	Surface         string         `json:"surface" gorm:"default:artificial" validate:"oneof=artificial grass synthetic"`
	HasLighting     bool           `json:"has_lighting" gorm:"default:false"`
	IsIndoor        bool           `json:"is_indoor" gorm:"default:false"`
	Amenities       []string       `json:"amenities" gorm:"type:json"`
	MaxPlayers      int            `json:"max_players" gorm:"default:4"`
	Rules           []string       `json:"rules" gorm:"type:json"`
	CancellationPolicy string      `json:"cancellation_policy"`
	AverageRating   float64        `json:"average_rating" gorm:"type:decimal(3,2);default:0"`
	ReviewCount     int            `json:"review_count" gorm:"default:0"`
	OwnerID         uint           `json:"owner_id" gorm:"not null"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Owner        User           `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Bookings     []Booking      `json:"bookings,omitempty" gorm:"foreignKey:CourtID"`
	Reviews      []Review       `json:"reviews,omitempty" gorm:"foreignKey:CourtID"`
	BusinessHours []BusinessHour `json:"business_hours,omitempty" gorm:"foreignKey:CourtID"`
	SpecialHours []SpecialHour  `json:"special_hours,omitempty" gorm:"foreignKey:CourtID"`
}

type BusinessHour struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CourtID   uint      `json:"court_id" gorm:"not null"`
	DayOfWeek int       `json:"day_of_week" gorm:"not null" validate:"min=0,max=6"` // 0=Sunday, 1=Monday, etc.
	OpenTime  string    `json:"open_time" gorm:"not null" validate:"required"`
	CloseTime string    `json:"close_time" gorm:"not null" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relación
	Court Court `json:"court,omitempty" gorm:"foreignKey:CourtID"`
}

type SpecialHour struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CourtID   uint      `json:"court_id" gorm:"not null"`
	Date      time.Time `json:"date" gorm:"type:date;not null"`
	OpenTime  *string   `json:"open_time,omitempty"`
	CloseTime *string   `json:"close_time,omitempty"`
	IsClosed  bool      `json:"is_closed" gorm:"default:false"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relación
	Court Court `json:"court,omitempty" gorm:"foreignKey:CourtID"`
}

type CreateCourtRequest struct {
	Name              string   `json:"name" validate:"required"`
	Address           string   `json:"address" validate:"required"`
	Latitude          float64  `json:"latitude" validate:"required"`
	Longitude         float64  `json:"longitude" validate:"required"`
	PricePerHour      float64  `json:"price_per_hour" validate:"required,min=0"`
	Description       string   `json:"description"`
	ImageURL          string   `json:"image_url"`
	Surface           string   `json:"surface" validate:"oneof=artificial grass synthetic"`
	HasLighting       bool     `json:"has_lighting"`
	IsIndoor          bool     `json:"is_indoor"`
	Amenities         []string `json:"amenities"`
	MaxPlayers        int      `json:"max_players" validate:"min=2,max=8"`
	Rules             []string `json:"rules"`
	CancellationPolicy string  `json:"cancellation_policy"`
	BusinessHours     []BusinessHourRequest `json:"business_hours" validate:"required,min=1"`
}

type BusinessHourRequest struct {
	DayOfWeek int    `json:"day_of_week" validate:"min=0,max=6"`
	OpenTime  string `json:"open_time" validate:"required"`
	CloseTime string `json:"close_time" validate:"required"`
}

type UpdateCourtRequest struct {
	Name              *string  `json:"name,omitempty"`
	Address           *string  `json:"address,omitempty"`
	Latitude          *float64 `json:"latitude,omitempty"`
	Longitude         *float64 `json:"longitude,omitempty"`
	PricePerHour      *float64 `json:"price_per_hour,omitempty" validate:"omitempty,min=0"`
	Description       *string  `json:"description,omitempty"`
	ImageURL          *string  `json:"image_url,omitempty"`
	Surface           *string  `json:"surface,omitempty" validate:"omitempty,oneof=artificial grass synthetic"`
	HasLighting       *bool    `json:"has_lighting,omitempty"`
	IsIndoor          *bool    `json:"is_indoor,omitempty"`
	Amenities         []string `json:"amenities,omitempty"`
	MaxPlayers        *int     `json:"max_players,omitempty" validate:"omitempty,min=2,max=8"`
	Rules             []string `json:"rules,omitempty"`
	CancellationPolicy *string `json:"cancellation_policy,omitempty"`
	IsActive          *bool    `json:"is_active,omitempty"`
}

type SearchCourtsRequest struct {
	MinPrice   *float64 `json:"min_price,omitempty" validate:"omitempty,min=0"`
	MaxPrice   *float64 `json:"max_price,omitempty" validate:"omitempty,min=0"`
	Surface    *string  `json:"surface,omitempty" validate:"omitempty,oneof=artificial grass synthetic"`
	HasLighting *bool   `json:"has_lighting,omitempty"`
	IsIndoor   *bool    `json:"is_indoor,omitempty"`
	MinRating  *float64 `json:"min_rating,omitempty" validate:"omitempty,min=0,max=5"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
	Radius     *float64 `json:"radius,omitempty" validate:"omitempty,min=0"`
}

type CourtStatistics struct {
	CourtID       uint    `json:"court_id"`
	CourtName     string  `json:"court_name"`
	TotalBookings int     `json:"total_bookings"`
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int     `json:"review_count"`
	Revenue       float64 `json:"revenue"`
}

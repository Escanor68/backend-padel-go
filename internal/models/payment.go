package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	BookingID     uint           `json:"booking_id" gorm:"not null"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	Amount        float64        `json:"amount" gorm:"type:decimal(10,2)" validate:"required,min=0"`
	Currency      string         `json:"currency" gorm:"default:ARS"`
	Status        string         `json:"status" gorm:"default:pending" validate:"oneof=pending approved rejected cancelled"`
	PreferenceID  string         `json:"preference_id" gorm:"uniqueIndex"`
	MercadoPagoID string         `json:"mercado_pago_id"`
	PaymentMethod string         `json:"payment_method"`
	PayerEmail    string         `json:"payer_email"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Booking Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	User    User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type CreatePreferenceRequest struct {
	BookingID  uint   `json:"booking_id" validate:"required"`
	PayerEmail string `json:"payer_email" validate:"required,email"`
}

type PreferenceResponse struct {
	ID        string `json:"id"`
	InitPoint string `json:"init_point"`
}

type PaymentStatusResponse struct {
	ID            string  `json:"id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
	PayerEmail    string  `json:"payer_email"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type WebhookRequest struct {
	Action string `json:"action"`
	Data   struct {
		ID string `json:"id"`
	} `json:"data"`
}

type PaymentStatistics struct {
	TotalPayments    int     `json:"total_payments"`
	ApprovedPayments int     `json:"approved_payments"`
	PendingPayments  int     `json:"pending_payments"`
	RejectedPayments int     `json:"rejected_payments"`
	TotalRevenue     float64 `json:"total_revenue"`
	AverageAmount    float64 `json:"average_amount"`
}

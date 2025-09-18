package services

import (
	"errors"
	"fmt"
	"time"

	"backend-padel-go/internal/config"
	"backend-padel-go/internal/models"

	"github.com/google/uuid"
	"github.com/mercadopago/sdk-go"
	"gorm.io/gorm"
)

type PaymentService struct {
	db *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{db: db}
}

func (s *PaymentService) CreatePreference(userID uint, req models.CreatePreferenceRequest) (*models.PreferenceResponse, error) {
	// Verificar que la reserva existe y pertenece al usuario
	var booking models.Booking
	if err := s.db.Where("id = ? AND user_id = ?", req.BookingID, userID).Preload("Court").First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, errors.New("failed to fetch booking")
	}

	// Verificar que la reserva no tenga ya un pago
	var existingPayment models.Payment
	if err := s.db.Where("booking_id = ?", req.BookingID).First(&existingPayment).Error; err == nil {
		return nil, errors.New("payment already exists for this booking")
	}

	// Crear pago en la base de datos
	paymentID := uuid.New().String()
	payment := models.Payment{
		ID:           paymentID,
		BookingID:    req.BookingID,
		UserID:       userID,
		Amount:       booking.TotalPrice,
		Currency:     "ARS",
		Status:       "pending",
		PayerEmail:   req.PayerEmail,
	}

	if err := s.db.Create(&payment).Error; err != nil {
		return nil, errors.New("failed to create payment")
	}

	// Crear preferencia en MercadoPago
	cfg := config.Load()
	client := mercadopago.NewClient(cfg.MercadoPago.AccessToken)

	preferenceRequest := mercadopago.PreferenceRequest{
		Items: []mercadopago.PreferenceItem{
			{
				Title:       fmt.Sprintf("Reserva - %s", booking.Court.Name),
				Description: fmt.Sprintf("Reserva para %s el %s de %s a %s", booking.Court.Name, booking.Date.Format("2006-01-02"), booking.StartTime, booking.EndTime),
				Quantity:    1,
				UnitPrice:   booking.TotalPrice,
			},
		},
		Payer: &mercadopago.PreferencePayer{
			Email: req.PayerEmail,
		},
		BackUrls: &mercadopago.PreferenceBackUrls{
			Success: fmt.Sprintf("%s/payment/success?payment_id=%s", cfg.MercadoPago.FrontendURL, paymentID),
			Failure: fmt.Sprintf("%s/payment/failure?payment_id=%s", cfg.MercadoPago.FrontendURL, paymentID),
			Pending: fmt.Sprintf("%s/payment/pending?payment_id=%s", cfg.MercadoPago.FrontendURL, paymentID),
		},
		AutoReturn: "approved",
		ExternalReference: paymentID,
		NotificationURL:   fmt.Sprintf("%s/api/v1/payments/webhook", cfg.MercadoPago.BackendURL),
	}

	preference, err := client.CreatePreference(preferenceRequest)
	if err != nil {
		// Eliminar pago de la base de datos si falla la creación de preferencia
		s.db.Delete(&payment)
		return nil, errors.New("failed to create MercadoPago preference")
	}

	// Actualizar pago con preference ID
	payment.PreferenceID = preference.ID
	if err := s.db.Save(&payment).Error; err != nil {
		return nil, errors.New("failed to update payment with preference ID")
	}

	return &models.PreferenceResponse{
		ID:        preference.ID,
		InitPoint: preference.InitPoint,
	}, nil
}

func (s *PaymentService) GetPaymentStatus(userID uint, paymentID string) (*models.PaymentStatusResponse, error) {
	var payment models.Payment
	if err := s.db.Where("id = ? AND user_id = ?", paymentID, userID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, errors.New("failed to fetch payment")
	}

	return &models.PaymentStatusResponse{
		ID:            payment.ID,
		Status:        payment.Status,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentMethod: payment.PaymentMethod,
		PayerEmail:    payment.PayerEmail,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}, nil
}

func (s *PaymentService) HandleWebhook(action string, data map[string]interface{}) error {
	paymentID, ok := data["id"].(string)
	if !ok {
		return errors.New("invalid payment ID in webhook data")
	}

	// Obtener información del pago desde MercadoPago
	cfg := config.Load()
	client := mercadopago.NewClient(cfg.MercadoPago.AccessToken)

	paymentInfo, err := client.GetPayment(paymentID)
	if err != nil {
		return errors.New("failed to get payment info from MercadoPago")
	}

	// Buscar pago en la base de datos por preference ID
	var payment models.Payment
	if err := s.db.Where("preference_id = ?", paymentInfo.PreferenceID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payment not found")
		}
		return errors.New("failed to fetch payment")
	}

	// Actualizar estado del pago
	payment.MercadoPagoID = paymentID
	payment.Status = s.mapMercadoPagoStatus(paymentInfo.Status)
	payment.PaymentMethod = paymentInfo.PaymentMethodID

	if err := s.db.Save(&payment).Error; err != nil {
		return errors.New("failed to update payment status")
	}

	// Si el pago fue aprobado, actualizar estado de la reserva
	if payment.Status == "approved" {
		if err := s.db.Model(&models.Booking{}).Where("id = ?", payment.BookingID).Update("status", "confirmed").Error; err != nil {
			return errors.New("failed to update booking status")
		}
	}

	return nil
}

func (s *PaymentService) mapMercadoPagoStatus(status string) string {
	switch status {
	case "approved":
		return "approved"
	case "pending":
		return "pending"
	case "rejected":
		return "rejected"
	case "cancelled":
		return "cancelled"
	default:
		return "pending"
	}
}

func (s *PaymentService) GetPaymentStatistics() (*models.PaymentStatistics, error) {
	var totalPayments int64
	var approvedPayments int64
	var pendingPayments int64
	var rejectedPayments int64
	var totalRevenue float64

	// Contar total de pagos
	if err := s.db.Model(&models.Payment{}).Count(&totalPayments).Error; err != nil {
		return nil, errors.New("failed to count total payments")
	}

	// Contar pagos por estado
	if err := s.db.Model(&models.Payment{}).Where("status = ?", "approved").Count(&approvedPayments).Error; err != nil {
		return nil, errors.New("failed to count approved payments")
	}

	if err := s.db.Model(&models.Payment{}).Where("status = ?", "pending").Count(&pendingPayments).Error; err != nil {
		return nil, errors.New("failed to count pending payments")
	}

	if err := s.db.Model(&models.Payment{}).Where("status = ?", "rejected").Count(&rejectedPayments).Error; err != nil {
		return nil, errors.New("failed to count rejected payments")
	}

	// Calcular ingresos totales
	if err := s.db.Model(&models.Payment{}).Where("status = ?", "approved").Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue).Error; err != nil {
		return nil, errors.New("failed to calculate total revenue")
	}

	// Calcular promedio
	var averageAmount float64
	if totalPayments > 0 {
		averageAmount = totalRevenue / float64(approvedPayments)
	}

	return &models.PaymentStatistics{
		TotalPayments:    int(totalPayments),
		ApprovedPayments: int(approvedPayments),
		PendingPayments:  int(pendingPayments),
		RejectedPayments: int(rejectedPayments),
		TotalRevenue:     totalRevenue,
		AverageAmount:    averageAmount,
	}, nil
}

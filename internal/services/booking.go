package services

import (
	"errors"
	"time"

	"backend-padel-go/internal/models"

	"gorm.io/gorm"
)

type BookingService struct {
	db *gorm.DB
}

func NewBookingService(db *gorm.DB) *BookingService {
	return &BookingService{db: db}
}

func (s *BookingService) CreateBooking(userID uint, req *models.CreateBookingRequest) (*models.BookingResponse, error) {
	// Verificar que la cancha existe
	var court models.Court
	if err := s.db.Where("id = ? AND is_active = ?", req.CourtID, true).First(&court).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("court not found")
		}
		return nil, errors.New("failed to fetch court")
	}

	// Parsear fecha
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	// Verificar disponibilidad
	available, err := s.checkAvailability(req.CourtID, date, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	if !available {
		return nil, errors.New("time slot not available")
	}

	// Calcular precio total
	hours := s.calculateHours(req.StartTime, req.EndTime)
	totalPrice := hours * court.PricePerHour

	// Crear reserva
	booking := models.Booking{
		CourtID:    req.CourtID,
		UserID:     userID,
		Date:       date,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Status:     "pending",
		TotalPrice: totalPrice,
		Notes:      req.Notes,
	}

	if err := s.db.Create(&booking).Error; err != nil {
		return nil, errors.New("failed to create booking")
	}

	// Cargar relaciones
	if err := s.db.Preload("Court").Preload("User").First(&booking, booking.ID).Error; err != nil {
		return nil, errors.New("failed to load booking with relations")
	}

	// Crear respuesta
	response := &models.BookingResponse{
		ID:         booking.ID,
		CourtID:    booking.CourtID,
		UserID:     booking.UserID,
		Date:       booking.Date,
		StartTime:  booking.StartTime,
		EndTime:    booking.EndTime,
		Status:     booking.Status,
		TotalPrice: booking.TotalPrice,
		Notes:      booking.Notes,
		CreatedAt:  booking.CreatedAt,
		UpdatedAt:  booking.UpdatedAt,
		Court: models.CourtInfo{
			ID:           booking.Court.ID,
			Name:         booking.Court.Name,
			Address:      booking.Court.Address,
			PricePerHour: booking.Court.PricePerHour,
			Surface:      booking.Court.Surface,
			HasLighting:  booking.Court.HasLighting,
			IsIndoor:     booking.Court.IsIndoor,
		},
		User: models.UserInfo{
			ID:        booking.User.ID,
			FirstName: booking.User.FirstName,
			LastName:  booking.User.LastName,
			Email:     booking.User.Email,
			Phone:     booking.User.Phone,
		},
	}

	return response, nil
}

func (s *BookingService) GetUserBookings(userID uint, filters *models.GetBookingsRequest) ([]*models.BookingResponse, error) {
	query := s.db.Model(&models.Booking{}).Where("user_id = ?", userID)

	// Aplicar filtros
	if filters.CourtID != nil {
		query = query.Where("court_id = ?", *filters.CourtID)
	}
	if filters.Date != nil {
		query = query.Where("date = ?", *filters.Date)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	var bookings []models.Booking
	if err := query.Preload("Court").Preload("User").Find(&bookings).Error; err != nil {
		return nil, errors.New("failed to fetch bookings")
	}

	// Convertir a respuesta
	var responses []*models.BookingResponse
	for _, booking := range bookings {
		responses = append(responses, &models.BookingResponse{
			ID:         booking.ID,
			CourtID:    booking.CourtID,
			UserID:     booking.UserID,
			Date:       booking.Date,
			StartTime:  booking.StartTime,
			EndTime:    booking.EndTime,
			Status:     booking.Status,
			TotalPrice: booking.TotalPrice,
			Notes:      booking.Notes,
			CreatedAt:  booking.CreatedAt,
			UpdatedAt:  booking.UpdatedAt,
			Court: models.CourtInfo{
				ID:           booking.Court.ID,
				Name:         booking.Court.Name,
				Address:      booking.Court.Address,
				PricePerHour: booking.Court.PricePerHour,
				Surface:      booking.Court.Surface,
				HasLighting:  booking.Court.HasLighting,
				IsIndoor:     booking.Court.IsIndoor,
			},
			User: models.UserInfo{
				ID:        booking.User.ID,
				FirstName: booking.User.FirstName,
				LastName:  booking.User.LastName,
				Email:     booking.User.Email,
				Phone:     booking.User.Phone,
			},
		})
	}

	return responses, nil
}

func (s *BookingService) GetBookingByID(id uint, userID uint) (*models.BookingResponse, error) {
	var booking models.Booking
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Preload("Court").Preload("User").First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, errors.New("failed to fetch booking")
	}

	response := &models.BookingResponse{
		ID:         booking.ID,
		CourtID:    booking.CourtID,
		UserID:     booking.UserID,
		Date:       booking.Date,
		StartTime:  booking.StartTime,
		EndTime:    booking.EndTime,
		Status:     booking.Status,
		TotalPrice: booking.TotalPrice,
		Notes:      booking.Notes,
		CreatedAt:  booking.CreatedAt,
		UpdatedAt:  booking.UpdatedAt,
		Court: models.CourtInfo{
			ID:           booking.Court.ID,
			Name:         booking.Court.Name,
			Address:      booking.Court.Address,
			PricePerHour: booking.Court.PricePerHour,
			Surface:      booking.Court.Surface,
			HasLighting:  booking.Court.HasLighting,
			IsIndoor:     booking.Court.IsIndoor,
		},
		User: models.UserInfo{
			ID:        booking.User.ID,
			FirstName: booking.User.FirstName,
			LastName:  booking.User.LastName,
			Email:     booking.User.Email,
			Phone:     booking.User.Phone,
		},
	}

	return response, nil
}

func (s *BookingService) CancelBooking(id uint, userID uint) error {
	var booking models.Booking
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return errors.New("failed to fetch booking")
	}

	// Verificar que la reserva se puede cancelar
	if booking.Status == "cancelled" {
		return errors.New("booking already cancelled")
	}

	if booking.Status == "completed" {
		return errors.New("cannot cancel completed booking")
	}

	// Actualizar estado
	if err := s.db.Model(&booking).Update("status", "cancelled").Error; err != nil {
		return errors.New("failed to cancel booking")
	}

	return nil
}

func (s *BookingService) checkAvailability(courtID uint, date time.Time, startTime, endTime string) (bool, error) {
	// Verificar si hay reservas existentes en el mismo horario
	var count int64
	err := s.db.Model(&models.Booking{}).
		Where("court_id = ? AND date = ? AND status = ?", courtID, date.Format("2006-01-02"), "confirmed").
		Where("(start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?)",
			endTime, startTime, endTime, startTime, startTime, endTime).
		Count(&count).Error

	if err != nil {
		return false, errors.New("failed to check availability")
	}

	return count == 0, nil
}

func (s *BookingService) calculateHours(startTime, endTime string) float64 {
	start, _ := time.Parse("15:04", startTime)
	end, _ := time.Parse("15:04", endTime)
	return end.Sub(start).Hours()
}

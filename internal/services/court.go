package services

import (
	"errors"
	"fmt"
	"time"

	"backend-padel-go/internal/models"

	"gorm.io/gorm"
)

type CourtService struct {
	db *gorm.DB
}

func NewCourtService(db *gorm.DB) *CourtService {
	return &CourtService{db: db}
}

func (s *CourtService) CreateCourt(ownerID uint, req models.CreateCourtRequest) (*models.Court, error) {
	// Crear cancha
	court := models.Court{
		Name:               req.Name,
		Address:            req.Address,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		PricePerHour:       req.PricePerHour,
		Description:        req.Description,
		ImageURL:           req.ImageURL,
		Surface:            req.Surface,
		HasLighting:        req.HasLighting,
		IsIndoor:           req.IsIndoor,
		Amenities:          req.Amenities,
		MaxPlayers:         req.MaxPlayers,
		Rules:              req.Rules,
		CancellationPolicy: req.CancellationPolicy,
		OwnerID:            ownerID,
		IsActive:           true,
	}

	if err := s.db.Create(&court).Error; err != nil {
		return nil, errors.New("failed to create court")
	}

	// Crear horarios de atención
	for _, bh := range req.BusinessHours {
		businessHour := models.BusinessHour{
			CourtID:   court.ID,
			DayOfWeek: bh.DayOfWeek,
			OpenTime:  bh.OpenTime,
			CloseTime: bh.CloseTime,
		}
		if err := s.db.Create(&businessHour).Error; err != nil {
			return nil, errors.New("failed to create business hours")
		}
	}

	// Cargar relaciones
	if err := s.db.Preload("BusinessHours").First(&court, court.ID).Error; err != nil {
		return nil, errors.New("failed to load court with business hours")
	}

	return &court, nil
}

func (s *CourtService) GetAllCourts() ([]models.Court, error) {
	var courts []models.Court
	if err := s.db.Preload("BusinessHours").Where("is_active = ?", true).Find(&courts).Error; err != nil {
		return nil, errors.New("failed to fetch courts")
	}
	return courts, nil
}

func (s *CourtService) GetCourtByID(id uint) (*models.Court, error) {
	var court models.Court
	if err := s.db.Preload("BusinessHours").Preload("Reviews").First(&court, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("court not found")
		}
		return nil, errors.New("failed to fetch court")
	}
	return &court, nil
}

func (s *CourtService) GetOwnerCourts(ownerID uint) ([]models.Court, error) {
	var courts []models.Court
	if err := s.db.Preload("BusinessHours").Where("owner_id = ?", ownerID).Find(&courts).Error; err != nil {
		return nil, errors.New("failed to fetch owner courts")
	}
	return courts, nil
}

func (s *CourtService) UpdateCourt(id uint, ownerID uint, req models.UpdateCourtRequest) (*models.Court, error) {
	var court models.Court
	if err := s.db.Where("id = ? AND owner_id = ?", id, ownerID).First(&court).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("court not found")
		}
		return nil, errors.New("failed to fetch court")
	}

	// Actualizar campos si se proporcionan
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.PricePerHour != nil {
		updates["price_per_hour"] = *req.PricePerHour
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.Surface != nil {
		updates["surface"] = *req.Surface
	}
	if req.HasLighting != nil {
		updates["has_lighting"] = *req.HasLighting
	}
	if req.IsIndoor != nil {
		updates["is_indoor"] = *req.IsIndoor
	}
	if req.Amenities != nil {
		updates["amenities"] = req.Amenities
	}
	if req.MaxPlayers != nil {
		updates["max_players"] = *req.MaxPlayers
	}
	if req.Rules != nil {
		updates["rules"] = req.Rules
	}
	if req.CancellationPolicy != nil {
		updates["cancellation_policy"] = *req.CancellationPolicy
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := s.db.Model(&court).Updates(updates).Error; err != nil {
		return nil, errors.New("failed to update court")
	}

	// Cargar la cancha actualizada
	if err := s.db.Preload("BusinessHours").First(&court, id).Error; err != nil {
		return nil, errors.New("failed to load updated court")
	}

	return &court, nil
}

func (s *CourtService) DeleteCourt(id uint, ownerID uint) error {
	var court models.Court
	if err := s.db.Where("id = ? AND owner_id = ?", id, ownerID).First(&court).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("court not found")
		}
		return errors.New("failed to fetch court")
	}

	if err := s.db.Delete(&court).Error; err != nil {
		return errors.New("failed to delete court")
	}

	return nil
}

func (s *CourtService) GetNearbyCourts(lat, lng, radius float64) ([]models.Court, error) {
	var courts []models.Court
	
	// Usar fórmula de Haversine para calcular distancia
	query := `
		SELECT *, 
		(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * 
		cos(radians(longitude) - radians(?)) + 
		sin(radians(?)) * sin(radians(latitude)))) AS distance
		FROM courts 
		WHERE is_active = true
		HAVING distance <= ?
		ORDER BY distance
	`

	if err := s.db.Raw(query, lat, lng, lat, radius).Scan(&courts).Error; err != nil {
		return nil, errors.New("failed to fetch nearby courts")
	}

	return courts, nil
}

func (s *CourtService) SearchCourts(req models.SearchCourtsRequest) ([]models.Court, error) {
	query := s.db.Model(&models.Court{}).Where("is_active = ?", true)

	if req.MinPrice != nil {
		query = query.Where("price_per_hour >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		query = query.Where("price_per_hour <= ?", *req.MaxPrice)
	}
	if req.Surface != nil {
		query = query.Where("surface = ?", *req.Surface)
	}
	if req.HasLighting != nil {
		query = query.Where("has_lighting = ?", *req.HasLighting)
	}
	if req.IsIndoor != nil {
		query = query.Where("is_indoor = ?", *req.IsIndoor)
	}
	if req.MinRating != nil {
		query = query.Where("average_rating >= ?", *req.MinRating)
	}

	// Si se proporcionan coordenadas, filtrar por distancia
	if req.Latitude != nil && req.Longitude != nil {
		radius := 20.0 // radio por defecto en km
		if req.Radius != nil {
			radius = *req.Radius
		}

		haversineQuery := `
			(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * 
			cos(radians(longitude) - radians(?)) + 
			sin(radians(?)) * sin(radians(latitude)))) <= ?
		`
		query = query.Where(haversineQuery, *req.Latitude, *req.Longitude, *req.Latitude, radius)
	}

	var courts []models.Court
	if err := query.Preload("BusinessHours").Find(&courts).Error; err != nil {
		return nil, errors.New("failed to search courts")
	}

	return courts, nil
}

func (s *CourtService) GetAvailability(courtID uint, date time.Time) ([]models.TimeSlot, error) {
	// Obtener cancha
	court, err := s.GetCourtByID(courtID)
	if err != nil {
		return nil, err
	}

	// Obtener horarios de atención para el día de la semana
	dayOfWeek := int(date.Weekday())
	var businessHour models.BusinessHour
	if err := s.db.Where("court_id = ? AND day_of_week = ?", courtID, dayOfWeek).First(&businessHour).Error; err != nil {
		return nil, errors.New("no business hours for this day")
	}

	// Obtener reservas existentes para la fecha
	var bookings []models.Booking
	if err := s.db.Where("court_id = ? AND date = ? AND status = ?", courtID, date.Format("2006-01-02"), "confirmed").Find(&bookings).Error; err != nil {
		return nil, errors.New("failed to fetch bookings")
	}

	// Generar slots de tiempo disponibles
	slots := s.generateTimeSlots(businessHour.OpenTime, businessHour.CloseTime, court.PricePerHour)
	
	// Filtrar slots ocupados
	availableSlots := s.filterAvailableSlots(slots, bookings)

	return availableSlots, nil
}

func (s *CourtService) generateTimeSlots(openTime, closeTime string, pricePerHour float64) []models.TimeSlot {
	var slots []models.TimeSlot
	
	// Parsear horarios
	open, _ := time.Parse("15:04", openTime)
	close, _ := time.Parse("15:04", closeTime)
	
	current := open
	for current.Before(close) {
		next := current.Add(time.Hour)
		if next.After(close) {
			break
		}
		
		slots = append(slots, models.TimeSlot{
			StartTime: current.Format("15:04"),
			EndTime:   next.Format("15:04"),
			Available: true,
			Price:     pricePerHour,
		})
		
		current = next
	}
	
	return slots
}

func (s *CourtService) filterAvailableSlots(slots []models.TimeSlot, bookings []models.Booking) []models.TimeSlot {
	availableSlots := make([]models.TimeSlot, 0)
	
	for _, slot := range slots {
		isAvailable := true
		for _, booking := range bookings {
			if s.isTimeSlotOverlapping(slot, booking) {
				isAvailable = false
				break
			}
		}
		
		if isAvailable {
			availableSlots = append(availableSlots, slot)
		}
	}
	
	return availableSlots
}

func (s *CourtService) isTimeSlotOverlapping(slot models.TimeSlot, booking models.Booking) bool {
	slotStart, _ := time.Parse("15:04", slot.StartTime)
	slotEnd, _ := time.Parse("15:04", slot.EndTime)
	bookingStart, _ := time.Parse("15:04", booking.StartTime)
	bookingEnd, _ := time.Parse("15:04", booking.EndTime)
	
	return (slotStart.Before(bookingEnd) && slotEnd.After(bookingStart))
}

func (s *CourtService) GetCourtStatistics(ownerID uint) ([]models.CourtStatistics, error) {
	var courts []models.Court
	if err := s.db.Preload("Bookings").Preload("Reviews").Where("owner_id = ?", ownerID).Find(&courts).Error; err != nil {
		return nil, errors.New("failed to fetch courts")
	}

	var statistics []models.CourtStatistics
	for _, court := range courts {
		totalBookings := len(court.Bookings)
		reviewCount := len(court.Reviews)
		
		var revenue float64
		for _, booking := range court.Bookings {
			revenue += booking.TotalPrice
		}
		
		statistics = append(statistics, models.CourtStatistics{
			CourtID:       court.ID,
			CourtName:     court.Name,
			TotalBookings: totalBookings,
			AverageRating: court.AverageRating,
			ReviewCount:   reviewCount,
			Revenue:       revenue,
		})
	}

	return statistics, nil
}

func (s *CourtService) CreateSpecialHours(courtID uint, ownerID uint, req models.SpecialHour) (*models.SpecialHour, error) {
	// Verificar que la cancha pertenece al propietario
	var court models.Court
	if err := s.db.Where("id = ? AND owner_id = ?", courtID, ownerID).First(&court).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("court not found")
		}
		return nil, errors.New("failed to fetch court")
	}

	specialHour := models.SpecialHour{
		CourtID:   courtID,
		Date:      req.Date,
		OpenTime:  req.OpenTime,
		CloseTime: req.CloseTime,
		IsClosed:  req.IsClosed,
		Reason:    req.Reason,
	}

	if err := s.db.Create(&specialHour).Error; err != nil {
		return nil, errors.New("failed to create special hours")
	}

	return &specialHour, nil
}

func (s *CourtService) GetSpecialHours(courtID uint, startDate, endDate time.Time) ([]models.SpecialHour, error) {
	var specialHours []models.SpecialHour
	if err := s.db.Where("court_id = ? AND date BETWEEN ? AND ?", courtID, startDate, endDate).Find(&specialHours).Error; err != nil {
		return nil, errors.New("failed to fetch special hours")
	}
	return specialHours, nil
}

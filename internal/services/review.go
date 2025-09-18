package services

import (
	"errors"

	"backend-padel-go/internal/models"

	"gorm.io/gorm"
)

type ReviewService struct {
	db *gorm.DB
}

func NewReviewService(db *gorm.DB) *ReviewService {
	return &ReviewService{db: db}
}

func (s *ReviewService) CreateReview(courtID uint, userID uint, req models.CreateReviewRequest) (*models.ReviewResponse, error) {
	// Verificar que la cancha existe
	var court models.Court
	if err := s.db.Where("id = ? AND is_active = ?", courtID, true).First(&court).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("court not found")
		}
		return nil, errors.New("failed to fetch court")
	}

	// Verificar que el usuario no haya ya reseñado esta cancha
	var existingReview models.Review
	if err := s.db.Where("court_id = ? AND user_id = ?", courtID, userID).First(&existingReview).Error; err == nil {
		return nil, errors.New("user has already reviewed this court")
	}

	// Crear reseña
	review := models.Review{
		CourtID: courtID,
		UserID:  userID,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	if err := s.db.Create(&review).Error; err != nil {
		return nil, errors.New("failed to create review")
	}

	// Actualizar estadísticas de la cancha
	if err := s.updateCourtRating(courtID); err != nil {
		return nil, errors.New("failed to update court rating")
	}

	// Cargar relaciones
	if err := s.db.Preload("Court").Preload("User").First(&review, review.ID).Error; err != nil {
		return nil, errors.New("failed to load review with relations")
	}

	// Crear respuesta
	response := &models.ReviewResponse{
		ID:        review.ID,
		CourtID:   review.CourtID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		User: models.UserInfo{
			ID:        review.User.ID,
			FirstName: review.User.FirstName,
			LastName:  review.User.LastName,
			Email:     review.User.Email,
		},
	}

	return response, nil
}

func (s *ReviewService) UpdateReview(reviewID uint, userID uint, req models.UpdateReviewRequest) (*models.ReviewResponse, error) {
	var review models.Review
	if err := s.db.Where("id = ? AND user_id = ?", reviewID, userID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, errors.New("failed to fetch review")
	}

	// Actualizar campos si se proporcionan
	updates := make(map[string]interface{})
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
	}

	if err := s.db.Model(&review).Updates(updates).Error; err != nil {
		return nil, errors.New("failed to update review")
	}

	// Actualizar estadísticas de la cancha
	if err := s.updateCourtRating(review.CourtID); err != nil {
		return nil, errors.New("failed to update court rating")
	}

	// Cargar relaciones
	if err := s.db.Preload("Court").Preload("User").First(&review, review.ID).Error; err != nil {
		return nil, errors.New("failed to load review with relations")
	}

	// Crear respuesta
	response := &models.ReviewResponse{
		ID:        review.ID,
		CourtID:   review.CourtID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		User: models.UserInfo{
			ID:        review.User.ID,
			FirstName: review.User.FirstName,
			LastName:  review.User.LastName,
			Email:     review.User.Email,
		},
	}

	return response, nil
}

func (s *ReviewService) DeleteReview(reviewID uint, userID uint) error {
	var review models.Review
	if err := s.db.Where("id = ? AND user_id = ?", reviewID, userID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("review not found")
		}
		return errors.New("failed to fetch review")
	}

	courtID := review.CourtID

	if err := s.db.Delete(&review).Error; err != nil {
		return errors.New("failed to delete review")
	}

	// Actualizar estadísticas de la cancha
	if err := s.updateCourtRating(courtID); err != nil {
		return errors.New("failed to update court rating")
	}

	return nil
}

func (s *ReviewService) GetCourtReviews(courtID uint) ([]models.ReviewResponse, error) {
	var reviews []models.Review
	if err := s.db.Where("court_id = ?", courtID).Preload("User").Order("created_at DESC").Find(&reviews).Error; err != nil {
		return nil, errors.New("failed to fetch reviews")
	}

	var responses []models.ReviewResponse
	for _, review := range reviews {
		responses = append(responses, models.ReviewResponse{
			ID:        review.ID,
			CourtID:   review.CourtID,
			UserID:    review.UserID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt,
			UpdatedAt: review.UpdatedAt,
			User: models.UserInfo{
				ID:        review.User.ID,
				FirstName: review.User.FirstName,
				LastName:  review.User.LastName,
				Email:     review.User.Email,
			},
		})
	}

	return responses, nil
}

func (s *ReviewService) updateCourtRating(courtID uint) error {
	var reviews []models.Review
	if err := s.db.Where("court_id = ?", courtID).Find(&reviews).Error; err != nil {
		return err
	}

	if len(reviews) == 0 {
		// Si no hay reseñas, establecer rating en 0
		return s.db.Model(&models.Court{}).Where("id = ?", courtID).Updates(map[string]interface{}{
			"average_rating": 0,
			"review_count":   0,
		}).Error
	}

	// Calcular promedio
	var totalRating float64
	for _, review := range reviews {
		totalRating += float64(review.Rating)
	}
	averageRating := totalRating / float64(len(reviews))

	// Actualizar cancha
	return s.db.Model(&models.Court{}).Where("id = ?", courtID).Updates(map[string]interface{}{
		"average_rating": averageRating,
		"review_count":   len(reviews),
	}).Error
}

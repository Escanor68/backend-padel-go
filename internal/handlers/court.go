package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend-padel-go/internal/middleware"
	"backend-padel-go/internal/models"
	"backend-padel-go/internal/services"

	"github.com/gin-gonic/gin"
)

type CourtHandler struct {
	courtService *services.CourtService
}

func NewCourtHandler(courtService *services.CourtService) *CourtHandler {
	return &CourtHandler{courtService: courtService}
}

// GetAllCourts godoc
// @Summary Get all courts
// @Description Get all active courts
// @Tags courts
// @Produce json
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /courts [get]
func (h *CourtHandler) GetAllCourts(c *gin.Context) {
	courts, err := h.courtService.GetAllCourts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch courts", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(courts))
}

// GetCourtByID godoc
// @Summary Get court by ID
// @Description Get a specific court by its ID
// @Tags courts
// @Produce json
// @Param id path int true "Court ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /courts/{id} [get]
func (h *CourtHandler) GetCourtByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid court ID", err.Error()))
		return
	}

	court, err := h.courtService.GetCourtByID(uint(id))
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Court not found", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch court", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(court))
}

// CreateCourt godoc
// @Summary Create a new court
// @Description Create a new court (owner only)
// @Tags courts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateCourtRequest true "Create court request"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /owner/courts [post]
func (h *CourtHandler) CreateCourt(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	ownerID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", "INTERNAL_ERROR"))
		return
	}

	var req models.CreateCourtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	court, err := h.courtService.CreateCourt(ownerID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Failed to create court", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(court))
}

// GetOwnerCourts godoc
// @Summary Get owner's courts
// @Description Get all courts owned by the authenticated user
// @Tags courts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /owner/courts [get]
func (h *CourtHandler) GetOwnerCourts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	ownerID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", "INTERNAL_ERROR"))
		return
	}

	courts, err := h.courtService.GetOwnerCourts(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch courts", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(courts))
}

// UpdateCourt godoc
// @Summary Update court
// @Description Update a court (owner only)
// @Tags courts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Court ID"
// @Param request body models.UpdateCourtRequest true "Update court request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /owner/courts/{id} [put]
func (h *CourtHandler) UpdateCourt(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	ownerID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", "INTERNAL_ERROR"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid court ID", err.Error()))
		return
	}

	var req models.UpdateCourtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	court, err := h.courtService.UpdateCourt(uint(id), ownerID, req)
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Court not found", err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Failed to update court", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(court))
}

// DeleteCourt godoc
// @Summary Delete court
// @Description Delete a court (owner only)
// @Tags courts
// @Produce json
// @Security BearerAuth
// @Param id path int true "Court ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /owner/courts/{id} [delete]
func (h *CourtHandler) DeleteCourt(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	ownerID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", "INTERNAL_ERROR"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid court ID", err.Error()))
		return
	}

	err = h.courtService.DeleteCourt(uint(id), ownerID)
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Court not found", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to delete court", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{"message": "Court deleted successfully"}))
}

// GetNearbyCourts godoc
// @Summary Get nearby courts
// @Description Get courts within a specified radius
// @Tags courts
// @Produce json
// @Param lat query number true "Latitude"
// @Param lng query number true "Longitude"
// @Param radius query number false "Radius in kilometers (default: 20)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /courts/nearby [get]
func (h *CourtHandler) GetNearbyCourts(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "20")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid latitude", err.Error()))
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid longitude", err.Error()))
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid radius", err.Error()))
		return
	}

	courts, err := h.courtService.GetNearbyCourts(lat, lng, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch nearby courts", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(courts))
}

// SearchCourts godoc
// @Summary Search courts
// @Description Search courts with filters
// @Tags courts
// @Produce json
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param surface query string false "Surface type"
// @Param has_lighting query boolean false "Has lighting"
// @Param is_indoor query boolean false "Is indoor"
// @Param min_rating query number false "Minimum rating"
// @Param lat query number false "Latitude"
// @Param lng query number false "Longitude"
// @Param radius query number false "Radius in kilometers"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /courts/search [get]
func (h *CourtHandler) SearchCourts(c *gin.Context) {
	var req models.SearchCourtsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid query parameters", err.Error()))
		return
	}

	courts, err := h.courtService.SearchCourts(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to search courts", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(courts))
}

// GetAvailability godoc
// @Summary Get court availability
// @Description Get available time slots for a court on a specific date
// @Tags courts
// @Produce json
// @Param id path int true "Court ID"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /courts/{id}/availability [get]
func (h *CourtHandler) GetAvailability(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid court ID", err.Error()))
		return
	}

	dateStr := c.Query("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid date format", err.Error()))
		return
	}

	slots, err := h.courtService.GetAvailability(uint(id), date)
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Court not found", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get availability", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(slots))
}

// GetCourtStatistics godoc
// @Summary Get court statistics
// @Description Get statistics for owner's courts
// @Tags courts
// @Produce json
// @Security BearerAuth
// @Param owner_id path int true "Owner ID"
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /owner/courts/{owner_id}/statistics [get]
func (h *CourtHandler) GetCourtStatistics(c *gin.Context) {
	ownerIDStr := c.Param("ownerId")
	ownerID, err := strconv.ParseUint(ownerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid owner ID", err.Error()))
		return
	}

	statistics, err := h.courtService.GetCourtStatistics(uint(ownerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get statistics", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(statistics))
}

// CreateSpecialHours godoc
// @Summary Create special hours
// @Description Create special hours for a court
// @Tags courts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Court ID"
// @Param request body models.SpecialHour true "Special hours request"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /owner/courts/{id}/special-hours [post]
func (h *CourtHandler) CreateSpecialHours(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	ownerID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", "INTERNAL_ERROR"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid court ID", err.Error()))
		return
	}

	var req models.SpecialHour
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	specialHour, err := h.courtService.CreateSpecialHours(uint(id), ownerID, req)
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Court not found", err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Failed to create special hours", err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(specialHour))
}

// GetSpecialHours godoc
// @Summary Get special hours
// @Description Get special hours for a court in a date range
// @Tags courts
// @Produce json
// @Param id path int true "Court ID"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /courts/{id}/special-hours [get]
func (h *CourtHandler) GetSpecialHours(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid court ID", err.Error()))
		return
	}

	startDateStr := c.Query("start_date")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid start date format", err.Error()))
		return
	}

	endDateStr := c.Query("end_date")
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid end date format", err.Error()))
		return
	}

	specialHours, err := h.courtService.GetSpecialHours(uint(id), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get special hours", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(specialHours))
}

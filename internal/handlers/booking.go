package handlers

import (
	"net/http"
	"strconv"

	"backend-padel-go/internal/middleware"
	"backend-padel-go/internal/models"
	"backend-padel-go/internal/services"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService *services.BookingService
}

func NewBookingHandler(bookingService *services.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking for a court
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateBookingRequest true "Create booking request"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", "INTERNAL_ERROR"))
		return
	}

	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	booking, err := h.bookingService.CreateBooking(userIDUint, req)
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Court not found", err.Error()))
		} else if err.Error() == "time slot not available" {
			c.JSON(http.StatusConflict, models.NewErrorResponse("Time slot not available", err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Failed to create booking", err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(booking))
}

// GetUserBookings godoc
// @Summary Get user bookings
// @Description Get all bookings for the authenticated user
// @Tags bookings
// @Produce json
// @Security BearerAuth
// @Param court_id query int false "Filter by court ID"
// @Param date query string false "Filter by date (YYYY-MM-DD)"
// @Param status query string false "Filter by status"
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /bookings [get]
func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", "INTERNAL_ERROR"))
		return
	}

	var filters models.GetBookingsRequest
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid query parameters", err.Error()))
		return
	}

	bookings, err := h.bookingService.GetUserBookings(userIDUint, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch bookings", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(bookings))
}

// GetBookingByID godoc
// @Summary Get booking by ID
// @Description Get a specific booking by its ID
// @Tags bookings
// @Produce json
// @Security BearerAuth
// @Param id path int true "Booking ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", "INTERNAL_ERROR"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	booking, err := h.bookingService.GetBookingByID(uint(id), userIDUint)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Booking not found", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch booking", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(booking))
}

// CancelBooking godoc
// @Summary Cancel booking
// @Description Cancel a booking
// @Tags bookings
// @Produce json
// @Security BearerAuth
// @Param id path int true "Booking ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /bookings/{id}/cancel [put]
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", "INTERNAL_ERROR"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	err = h.bookingService.CancelBooking(uint(id), userIDUint)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Booking not found", err.Error()))
		} else if err.Error() == "booking already cancelled" || err.Error() == "cannot cancel completed booking" {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Cannot cancel booking", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to cancel booking", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{"message": "Booking cancelled successfully"}))
}

package handlers

import (
	"net/http"
	"strconv"

	"backend-padel-go/internal/models"
	"backend-padel-go/internal/services"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *services.ReviewService
}

func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// CreateReview godoc
// @Summary Create a review
// @Description Create a review for a court
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Court ID"
// @Param request body models.CreateReviewRequest true "Create review request"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /reviews/courts/{id} [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
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
	courtID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid court ID", err.Error()))
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	review, err := h.reviewService.CreateReview(uint(courtID), userIDUint, req)
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Court not found", err.Error()))
		} else if err.Error() == "user has already reviewed this court" {
			c.JSON(http.StatusConflict, models.NewErrorResponse("User has already reviewed this court", err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Failed to create review", err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(review))
}

// UpdateReview godoc
// @Summary Update review
// @Description Update a review
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Review ID"
// @Param request body models.UpdateReviewRequest true "Update review request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
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
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid review ID", err.Error()))
		return
	}

	var req models.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	review, err := h.reviewService.UpdateReview(uint(reviewID), userIDUint, req)
	if err != nil {
		if err.Error() == "review not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Review not found", err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Failed to update review", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(review))
}

// DeleteReview godoc
// @Summary Delete review
// @Description Delete a review
// @Tags reviews
// @Produce json
// @Security BearerAuth
// @Param id path int true "Review ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
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
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid review ID", err.Error()))
		return
	}

	err = h.reviewService.DeleteReview(uint(reviewID), userIDUint)
	if err != nil {
		if err.Error() == "review not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Review not found", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to delete review", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{"message": "Review deleted successfully"}))
}

// GetCourtReviews godoc
// @Summary Get court reviews
// @Description Get all reviews for a specific court
// @Tags reviews
// @Produce json
// @Param id path int true "Court ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /courts/{id}/reviews [get]
func (h *ReviewHandler) GetCourtReviews(c *gin.Context) {
	idStr := c.Param("id")
	courtID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid court ID", err.Error()))
		return
	}

	reviews, err := h.reviewService.GetCourtReviews(uint(courtID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch reviews", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(reviews))
}

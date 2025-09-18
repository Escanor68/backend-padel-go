package handlers

import (
	"net/http"
	"strconv"

	"backend-padel-go/internal/middleware"
	"backend-padel-go/internal/models"
	"backend-padel-go/internal/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreatePreference godoc
// @Summary Create payment preference
// @Description Create a MercadoPago payment preference
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreatePreferenceRequest true "Create preference request"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /payments/preference [post]
func (h *PaymentHandler) CreatePreference(c *gin.Context) {
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

	var req models.CreatePreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	preference, err := h.paymentService.CreatePreference(userIDUint, req)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Booking not found", err.Error()))
		} else if err.Error() == "payment already exists for this booking" {
			c.JSON(http.StatusConflict, models.NewErrorResponse("Payment already exists for this booking", err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Failed to create preference", err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(preference))
}

// GetPaymentStatus godoc
// @Summary Get payment status
// @Description Get the status of a payment
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /payments/{id}/status [get]
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
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

	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Payment ID is required", "BAD_REQUEST"))
		return
	}

	status, err := h.paymentService.GetPaymentStatus(userIDUint, paymentID)
	if err != nil {
		if err.Error() == "payment not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Payment not found", err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get payment status", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(status))
}

// HandleWebhook godoc
// @Summary Handle MercadoPago webhook
// @Description Handle webhook notifications from MercadoPago
// @Tags payments
// @Accept json
// @Produce json
// @Param request body models.WebhookRequest true "Webhook request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /payments/webhook [post]
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	var req models.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Convertir data a map[string]interface{}
	data := make(map[string]interface{})
	data["id"] = req.Data.ID

	err := h.paymentService.HandleWebhook(req.Action, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to handle webhook", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{"received": true}))
}

// GetPaymentStatistics godoc
// @Summary Get payment statistics
// @Description Get payment statistics (admin only)
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /payments/statistics [get]
func (h *PaymentHandler) GetPaymentStatistics(c *gin.Context) {
	// Verificar que el usuario es admin
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", "UNAUTHORIZED"))
		return
	}

	roleStr, ok := userRole.(string)
	if !ok || roleStr != "admin" {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Admin access required", "FORBIDDEN"))
		return
	}

	statistics, err := h.paymentService.GetPaymentStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get payment statistics", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(statistics))
}

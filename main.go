package main

import (
	"backend-padel-go/internal/config"
	"backend-padel-go/internal/database"
	"backend-padel-go/internal/handlers"
	"backend-padel-go/internal/middleware"
	"backend-padel-go/internal/models"
	"backend-padel-go/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "backend-padel-go/docs" // Importar documentación Swagger
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Padel Courts API
// @version 1.0
// @description API para gestión de canchas de padel
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Configurar base de datos
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrar esquemas
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Middleware global
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Inicializar servicios
	authService := services.NewAuthService(db)
	courtService := services.NewCourtService(db)
	bookingService := services.NewBookingService(db)
	paymentService := services.NewPaymentService(db)
	reviewService := services.NewReviewService(db)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	courtHandler := handlers.NewCourtHandler(courtService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	// Configurar Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configurar rutas
	setupRoutes(r, authHandler, courtHandler, bookingHandler, paymentHandler, reviewHandler)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(r *gin.Engine, authHandler *handlers.AuthHandler, courtHandler *handlers.CourtHandler, bookingHandler *handlers.BookingHandler, paymentHandler *handlers.PaymentHandler, reviewHandler *handlers.ReviewHandler) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Grupo de rutas API
	v1 := r.Group("/api/v1")
	{
		// Rutas de autenticación
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Rutas públicas de canchas
		courts := v1.Group("/courts")
		{
			courts.GET("", courtHandler.GetAllCourts)
			courts.GET("/:id", courtHandler.GetCourtByID)
			courts.GET("/nearby", courtHandler.GetNearbyCourts)
			courts.GET("/search", courtHandler.SearchCourts)
			courts.GET("/:id/availability", courtHandler.GetAvailability)
			courts.GET("/:id/reviews", reviewHandler.GetCourtReviews)
		}

		// Rutas protegidas
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// Gestión de canchas (propietarios)
			owner := protected.Group("/owner")
			{
				owner.POST("/courts", courtHandler.CreateCourt)
				owner.GET("/courts", courtHandler.GetOwnerCourts)
				owner.PUT("/courts/:id", courtHandler.UpdateCourt)
				owner.DELETE("/courts/:id", courtHandler.DeleteCourt)
				owner.GET("/courts/:id/statistics", courtHandler.GetCourtStatistics)
				owner.POST("/courts/:id/special-hours", courtHandler.CreateSpecialHours)
				owner.GET("/courts/:id/special-hours", courtHandler.GetSpecialHours)
			}

			// Reservas
			bookings := protected.Group("/bookings")
			{
				bookings.POST("", bookingHandler.CreateBooking)
				bookings.GET("", bookingHandler.GetUserBookings)
				bookings.GET("/:id", bookingHandler.GetBookingByID)
				bookings.PUT("/:id/cancel", bookingHandler.CancelBooking)
			}

			// Reseñas
			reviews := protected.Group("/reviews")
			{
				reviews.POST("/courts/:id", reviewHandler.CreateReview)
				reviews.PUT("/:id", reviewHandler.UpdateReview)
				reviews.DELETE("/:id", reviewHandler.DeleteReview)
			}

			// Pagos
			payments := protected.Group("/payments")
			{
				payments.POST("/preference", paymentHandler.CreatePreference)
				payments.GET("/:id/status", paymentHandler.GetPaymentStatus)
				payments.POST("/webhook", paymentHandler.HandleWebhook)
			}
		}
	}
}

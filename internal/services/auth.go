package services

import (
	"errors"
	"time"

	"backend-padel-go/internal/config"
	"backend-padel-go/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Verificar si el usuario ya existe
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Crear usuario
	user := models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      "user",
		IsActive:  true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generar tokens
	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.generateRefreshToken(&user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	cfg := config.Load()
	expiresIn := int64(cfg.JWT.AccessTokenExpiry * 3600) // convertir horas a segundos

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generar tokens
	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.generateRefreshToken(&user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	cfg := config.Load()
	expiresIn := int64(cfg.JWT.AccessTokenExpiry * 3600) // convertir horas a segundos

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	// Validar refresh token
	claims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Obtener usuario
	var user models.User
	if err := s.db.Where("id = ? AND is_active = ?", claims.UserID, true).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Generar nuevos tokens
	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	newRefreshToken, err := s.generateRefreshToken(&user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	cfg := config.Load()
	expiresIn := int64(cfg.JWT.AccessTokenExpiry * 3600) // convertir horas a segundos

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	cfg := config.Load()
	
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.AccessTokenExpiry) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.SecretKey))
}

func (s *AuthService) generateRefreshToken(user *models.User) (string, error) {
	cfg := config.Load()
	
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.RefreshTokenExpiry) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.SecretKey))
}

func (s *AuthService) validateRefreshToken(tokenString string) (*Claims, error) {
	cfg := config.Load()
	
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ValidateToken(tokenString string) (*Claims, error) {
	cfg := config.Load()
	
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

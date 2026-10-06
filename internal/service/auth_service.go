package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/saddam/employee-management-be/internal/config"
	"github.com/saddam/employee-management-be/internal/middleware"
	"github.com/saddam/employee-management-be/internal/model"
	"github.com/saddam/employee-management-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(name, email, password string) (*LoginResponse, error)
	Login(email, password string) (*LoginResponse, error)
	GetMe(userID uint) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(name, email, password string) (*LoginResponse, error) {
	existing, err := s.userRepo.FindByEmail(email)
	if err == nil && existing != nil {
		return nil, errors.New("email is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     model.RoleViewer,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, errors.New("failed to create user account")
	}

	expirationTime := time.Now().Add(time.Duration(s.cfg.JWTExpireHours) * time.Hour)
	claims := &middleware.JWTClaims{
		UserID: newUser.ID,
		Email:  newUser.Email,
		Role:   newUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResponse{
		Token: tokenString,
		User:  newUser,
	}, nil
}

func (s *authService) Login(email, password string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	expirationTime := time.Now().Add(time.Duration(s.cfg.JWTExpireHours) * time.Hour)
	claims := &middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

func (s *authService) GetMe(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

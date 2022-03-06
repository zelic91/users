package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"zelic91/users/apple"
	"zelic91/users/gen/models"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo Repo
}

const (
	ErrInvalidUsername  = "your username is already existed"
	ErrPasswordTooShort = "your password is too short"
	ErrPasswordNotMatch = "your password does not match"
	ErrInvalidPassword  = "your password is invalid"

	MinPasswordLength = 6
	HashCost          = 2

	JWTSecret = "Hello App"
)

func (s AuthService) SignUp(ctx context.Context, params *models.SignUpRequest) (*models.AuthResponse, error) {
	password := params.Password

	if len(*password) < MinPasswordLength {
		return nil, errors.New(ErrPasswordTooShort)
	}

	if strings.Compare(*password, *params.PasswordConfirmation) != 0 {
		return nil, errors.New(ErrPasswordNotMatch)
	}

	hashedPassword, err := hashPassword(*password)
	if err != nil {
		return nil, err
	}

	user := User{
		Username:       *params.Username,
		HashedPassword: hashedPassword,
	}

	output, err := s.Repo.Create(&user)
	if err != nil {
		return nil, err
	}

	accessToken, err := generateToken(output)
	if err != nil {
		return nil, err
	}

	return toAuthResponse(output, *accessToken), nil
}

func (s AuthService) SignIn(ctx context.Context, params *models.SignInRequest) (*models.AuthResponse, error) {
	user, err := s.Repo.GetByUsername(*params.Username)
	if err != nil {
		return nil, err
	}

	if !verifyPassword(*params.Password, user.HashedPassword) {
		return nil, errors.New(ErrInvalidPassword)
	}

	accessToken, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	return toAuthResponse(user, *accessToken), nil
}

func (s AuthService) SignInWithApple(ctx context.Context, params *models.SignInAppleRequest) (*models.AuthResponse, error) {
	// Verify token
	tokenString := params.Token

	// Parse info
	token, err := jwt.ParseWithClaims(*tokenString, &apple.AppleClaims{}, apple.AuthFunc)

	if err != nil {
		return nil, err
	}

	// Check if user exists
	claims := token.Claims.(*apple.AppleClaims)
	email := claims.Email

	user, err := s.Repo.GetByUsername(email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Create new user
	if user == nil {
		password := strings.Repeat(claims.Email, 3)
		hashedPassword, err := hashPassword(password)
		if err != nil {
			return nil, err
		}

		newUser := User{
			Username:       claims.Email,
			Email:          &claims.Email,
			HashedPassword: hashedPassword,
		}

		user, err = s.Repo.Create(&newUser)
		if err != nil {
			return nil, err
		}
	}

	accessToken, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	return toAuthResponse(user, *accessToken), nil
}

func generateToken(user *User) (*string, error) {
	key := []byte(JWTSecret)
	claim := UserClaims{
		ID:       user.ID,
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedString, err := token.SignedString(key)

	if err != nil {
		return nil, err
	}

	return &signedString, nil
}

func hashPassword(rawPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), HashCost)
	return string(bytes), err
}

func verifyPassword(rawPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	return err == nil
}

func toAuthResponse(user *User, accessToken string) *models.AuthResponse {
	return &models.AuthResponse{
		AccessToken: accessToken,
		Username:    user.Username,
	}
}

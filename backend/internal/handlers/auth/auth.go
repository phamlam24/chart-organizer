package auth

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"log"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"

	authv1 "chart-organizer/backend/gen/contracts/auth/v1"
	"chart-organizer/backend/internal/interceptors"
	authRepo "chart-organizer/backend/internal/repository/auth"
)

type AuthHandler struct {
	DB *sql.DB
}

func generateJWT(username string, userID string) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	currentTime := time.Now()
	expirationTime := currentTime.Add(24 * time.Hour)
	claims := &interceptors.Claims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthHandler) Signup(
	ctx context.Context,
	req *connect.Request[authv1.SignupRequest],
) (*connect.Response[authv1.SignupResponse], error) {
	log.Println("Request headers: ", req.Header())

	username := req.Msg.Username
	password := req.Msg.Password

	// Validate input
	if username == "" || password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username and password are required"))
	}

	// Try to add the user
	err := authRepo.AddNewUser(s.DB, username, password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to create user: "+err.Error()))
	}

	// Get the user ID
	userID, err := authRepo.GetUserID(s.DB, username)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to retrieve user ID"))
	}

	// Generate JWT
	token, err := generateJWT(username, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate token"))
	}

	res := connect.NewResponse(&authv1.SignupResponse{
		JwtToken: token,
	})
	return res, nil
}

func (s *AuthHandler) Login(
	ctx context.Context,
	req *connect.Request[authv1.LoginRequest],
) (*connect.Response[authv1.LoginResponse], error) {
	log.Println("Request headers: ", req.Header())

	username := req.Msg.Username
	password := req.Msg.Password

	// Validate input
	if username == "" || password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username and password are required"))
	}

	// Check credentials
	isValid, err := authRepo.CheckUsernameAndPassword(s.DB, username, password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("authentication failed: "+err.Error()))
	}

	if !isValid {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid username or password"))
	}

	// Get the user ID
	userID, err := authRepo.GetUserID(s.DB, username)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to retrieve user ID"))
	}

	// Generate JWT
	token, err := generateJWT(username, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate token"))
	}

	res := connect.NewResponse(&authv1.LoginResponse{
		JwtToken: token,
	})
	return res, nil
}

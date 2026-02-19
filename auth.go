// ============================================================
// AUTH SERVICE — internal/services/auth.go
// ============================================================
// Business logic for authentication.
// Handlers call this. This calls the repository (DB).
//
// GO CONCEPTS TO LEARN:
//   - interfaces: AuthService is an interface.
//     Any struct that has Signup(), Login(), RefreshToken() methods
//     automatically implements it. No "implements" keyword needed.
//   - context.Context: pass this through every function that touches DB/Redis
//     It carries deadlines and cancellation signals.
//   - error wrapping: fmt.Errorf("signup failed: %w", err)
//     %w lets you unwrap the error later
//
// IMPLEMENT IN THIS ORDER:
//   1. Signup()
//   2. Login()
//   3. RefreshToken()
// ============================================================

package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/movie-reservation/internal/models"
	"github.com/yourusername/movie-reservation/internal/repository"
)

// Sentinel errors — return these so handlers can check the error type
// LEARN: Sentinel errors are like named constants for errors
var (
	ErrEmailTaken      = errors.New("email already taken")
	ErrInvalidCreds    = errors.New("invalid credentials")
	ErrInvalidToken    = errors.New("invalid token")
)

// AuthService interface — the contract this service must fulfill
// LEARN: Why an interface? So you can swap implementations and write tests easily.
type AuthService interface {
	Signup(ctx context.Context, req models.SignupRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

// authService is the concrete implementation
type authService struct {
	userRepo repository.UserRepository // talks to DB
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// ── SIGNUP ───────────────────────────────────────────────────
// TODO (Phase 1 - Step 1): Implement this first
func (s *authService) Signup(ctx context.Context, req models.SignupRequest) (*models.AuthResponse, error) {
	// Step 1: Check if email already exists
	// TODO: existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	// TODO: if existing != nil {
	//   return nil, ErrEmailTaken
	// }

	// Step 2: Hash the password
	// LEARN: bcrypt is slow BY DESIGN. Makes brute force attacks expensive.
	// Cost 12 = ~300ms per hash. That's fine for login, terrible for brute force.
	// TODO: hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	// TODO: if err != nil {
	//   return nil, fmt.Errorf("hashing password: %w", err)
	// }

	// Step 3: Save user to DB
	// TODO: user, err := s.userRepo.Create(ctx, repository.CreateUserParams{
	//   Email:        req.Email,
	//   PasswordHash: string(hash),
	//   Role:         "USER",
	// })
	// TODO: if err != nil {
	//   return nil, fmt.Errorf("creating user: %w", err)
	// }

	// Step 4: Generate tokens
	// TODO: accessToken, err := generateAccessToken(user.ID, user.Role)
	// TODO: refreshToken, err := generateRefreshToken(user.ID)

	// TODO: return &models.AuthResponse{
	//   AccessToken:  accessToken,
	//   RefreshToken: refreshToken,
	// }, nil

	return nil, errors.New("TODO: implement signup")
}

// ── LOGIN ────────────────────────────────────────────────────
// TODO (Phase 1 - Step 2)
func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	// Step 1: Find user by email
	// TODO: user, err := s.userRepo.FindByEmail(ctx, req.Email)
	// TODO: if err != nil || user == nil {
	//   return nil, ErrInvalidCreds // don't reveal if email exists
	// }

	// Step 2: Compare password
	// bcrypt.CompareHashAndPassword returns nil if match, error if not
	// TODO: if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
	//   return nil, ErrInvalidCreds
	// }

	// Step 3: Generate tokens
	// TODO: accessToken, _ := generateAccessToken(user.ID, user.Role)
	// TODO: refreshToken, _ := generateRefreshToken(user.ID)
	// TODO: return &models.AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil

	return nil, errors.New("TODO: implement login")
}

// ── REFRESH TOKEN ────────────────────────────────────────────
// TODO (Phase 1 - Step 3)
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// TODO: Parse and verify the refresh token
	// TODO: Extract userID from claims
	// TODO: Generate and return a new access token

	return "", errors.New("TODO: implement refresh token")
}

// ── JWT HELPERS ───────────────────────────────────────────────
// LEARN: A JWT has 3 parts: header.payload.signature
//   Header: algorithm used
//   Payload: your data (claims) — userId, role, expiry
//   Signature: proves it wasn't tampered with
//
// Access token: short-lived (15 min), used for every request
// Refresh token: long-lived (7 days), only used to get new access tokens

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func generateAccessToken(userID, role string) (string, error) {
	// TODO: Implement this
	// claims := Claims{
	//   UserID: userID,
	//   Role:   role,
	//   RegisteredClaims: jwt.RegisteredClaims{
	//     ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	//     IssuedAt:  jwt.NewNumericDate(time.Now()),
	//   },
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	return "", errors.New("TODO")
}

func generateRefreshToken(userID string) (string, error) {
	// TODO: Similar to above but longer expiry, no role claim
	// ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	return "", errors.New("TODO")
}

func parseToken(tokenStr, secret string) (*Claims, error) {
	// TODO: Parse and validate a JWT string
	// token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
	//   if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
	//     return nil, fmt.Errorf("unexpected signing method")
	//   }
	//   return []byte(secret), nil
	// })
	// if err != nil || !token.Valid {
	//   return nil, ErrInvalidToken
	// }
	// return token.Claims.(*Claims), nil
	return nil, errors.New("TODO")
}

// Keep compiler happy until implemented
var _ = bcrypt.DefaultCost
var _ = fmt.Sprintf
var _ = time.Now
var _ = os.Getenv

package handlers

import (
	"encoding/json"
	"net/http"

	"kafka-blog-backend/internal/models"
	"kafka-blog-backend/internal/services"

	"github.com/go-chi/jwtauth/v5"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get token and claims from context (set by jwtauth.Verifier)
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Token error", http.StatusUnauthorized)
		return
	}

	// Extract user info from claims
	userID, _ := claims["user_id"].(float64) // JSON numbers are float64
	username, _ := claims["username"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	// TODO: Get full user details from database
	// user, err := h.userService.GetUserByID(int(userID))

	// For now, return user from claims
	user := models.User{
		ID:       int(userID),
		Username: username,
		Email:    email,
		Role:     role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Validate request (you might want to use a validation library)

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// TODO: Save user to database
	// user, err := h.userService.CreateUser(req.Username, req.Email, hashedPassword)
	// For now, creating a mock user response
	user := models.User{
		ID:       1, // This should come from database
		Username: req.Username,
		Email:    req.Email,
		Role:     "user", // Default role
		Password: hashedPassword,
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Get user from database by email
	// user, err := h.userService.GetUserByEmail(req.Email)
	// For now, creating a mock user with a test password
	// The password "password123" hashed with bcrypt
	testPasswordHash := "$2a$10$YourHashedPasswordHere" // You'll need to generate this

	user := models.User{
		ID:       1,
		Username: "testuser",
		Email:    req.Email,
		Password: testPasswordHash,
		Role:     "user",
	}

	// Check password
	if err := h.authService.CheckPassword(req.Password, user.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Don't include password in response
	user.Password = ""

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

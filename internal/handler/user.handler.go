package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NekKkMirror/erm-tz/internal/dto"
	"github.com/NekKkMirror/erm-tz/internal/service"
	"github.com/gorilla/mux"
)

// UserHandler handles HTTP requests related to user operations such as registration and email verification.
// It interacts with the UserService and EmailService to perform these operations.
type UserHandler struct {
	userService  *service.UserService
	emailService *service.EmailService
}

// NewUserHandler initializes and returns a new UserHandler instance with the provided UserService and EmailService.
func NewUserHandler(userService *service.UserService, emailService *service.EmailService) *UserHandler {
	return &UserHandler{userService: userService, emailService: emailService}
}

// Register handles the registration of new users.
// It decodes the incoming JSON request, interacts with UserService to create a new user,
// and sends a verification email with the help of EmailService.
func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	err := uh.userService.Register(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	resp := dto.UserResponse{
		Message: fmt.Sprintf("A verification email has been sent to %s. Please verify your email within the next %d minutes to complete your registration.", req.Email, uh.emailService.GetTokenExpiry()),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// VerifyEmail handles the email verification process.
// It extracts the token from the request, verifies it using the UserService, and updates the user's verification status.
func (uh *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	err := uh.userService.VerifyEmail(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.UserResponse{Message: "Email successfully verified"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RegisterUsersRouter registers the HTTP routes for the user endpoints.
// It sets up the routes for user registration and email verification.
func RegisterUsersRouter(router *mux.Router, handler *UserHandler) {
	router.HandleFunc("/users/register", handler.Register).Methods("POST")
	router.HandleFunc("/users/verify", handler.VerifyEmail).Methods("GET")
}

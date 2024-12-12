package service

import (
	"errors"
	"time"

	"github.com/NekKkMirror/erm-tz/internal/dto"
	"github.com/NekKkMirror/erm-tz/internal/model"
	"github.com/NekKkMirror/erm-tz/internal/repository"
	"github.com/NekKkMirror/erm-tz/internal/utils"
)

// UserService is responsible for managing user-related operations,
// such as registration, email verification, and database interactions.
type UserService struct {
	repo         *repository.UserRepository
	emailService *EmailService
}

// NewUserService creates a new UserService instance with the provided UserRepository and EmailService.
func NewUserService(repo *repository.UserRepository, emailService *EmailService) *UserService {
	return &UserService{repo: repo, emailService: emailService}
}

// Register creates a new user record in the database with the provided registration data,
// generates a verification token, and sends a verification email to the user.
func (us *UserService) Register(req dto.UserRegisterRequest) error {
	u := &model.User{
		Nickname:  req.Nickname,
		Email:     req.Email,
		Verified:  false,
		CreatedAt: time.Now(),
	}

	err := us.repo.Save(u)
	if err != nil {
		return err
	}

	token, err := us.emailService.GenerateVerificationToken(req.Email)
	if err != nil {
		return err
	}

	return us.emailService.SendVerificationEmail(req.Email, token)
}

// VerifyEmail verifies the user's email using the provided token by decoding the token,
// extracting the email, and updating the user's verification status in the database.
func (us *UserService) VerifyEmail(token string) error {
	claims, err := utils.DecodeJWT(token, us.emailService.JWTSecretKey)
	if err != nil {
		return err
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return errors.New("invalid token data")
	}

	err = us.repo.VerifyEmail(email)
	if err != nil {
		return err
	}

	return nil
}

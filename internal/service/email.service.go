package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/NekKkMirror/erm-tz/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// EmailService manages email-related functionality including generating JWT-based verification tokens,
// refreshing OAuth2 tokens, and sending emails using the Gmail API.
type EmailService struct {
	JWTSecretKey    string
	VerificationURL string
	TokenExpiry     int

	GoogleOAuthAccessToken  string
	GoogleOAuthRefreshToken string
	oauthConfig             *oauth2.Config
}

// NewEmailService initializes and returns a new EmailService instance using the provided configuration.
func NewEmailService(config *config.Config) *EmailService {
	oauthConfig := &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  config.GoogleRedirectURI,
		Scopes:       []string{gmail.GmailSendScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.GoogleAuthURL,
			TokenURL: config.GoogleTokenURL,
		},
	}

	return &EmailService{
		JWTSecretKey:    config.JWTSecretKey,
		VerificationURL: config.EmailVerificationURL,
		TokenExpiry:     config.EmailTokenExpiry,

		GoogleOAuthAccessToken:  config.GoogleOAuthAccessToken,
		GoogleOAuthRefreshToken: config.GoogleOAuthRefreshToken,
		oauthConfig:             oauthConfig,
	}
}

// GenerateVerificationToken creates a JWT token for a given email address with a configured expiration time.
// The token can be used for email verification purposes.
func (es *EmailService) GenerateVerificationToken(email string) (string, error) {
	exp := time.Now().Add(time.Minute * time.Duration(es.TokenExpiry)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   exp,
	})

	tokenString, err := token.SignedString([]byte(es.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}

// SendVerificationEmail sends an email containing a verification link to the given email address
// using the Gmail API. It refreshes the OAuth2 token if necessary.
func (es *EmailService) SendVerificationEmail(email, token string) error {
	oauthToken := &oauth2.Token{
		AccessToken:  es.GoogleOAuthAccessToken,
		RefreshToken: es.GoogleOAuthRefreshToken,
	}

	tokenSource := es.oauthConfig.TokenSource(context.Background(), oauthToken)

	newToken, err := tokenSource.Token()
	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	// log someone
	log.Printf("Access token will expire at: %v", newToken.Expiry)

	service, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Failed to create Gmail service: %v", err)
		return fmt.Errorf("failed to create Gmail service: %w", err)
	}

	messageBody := fmt.Sprintf(
		"To verify your account, please click the link below within %d minutes:\n%s?token=%s",
		es.TokenExpiry, es.VerificationURL, token,
	)

	var message gmail.Message
	message.Raw = base64.URLEncoding.EncodeToString([]byte(
		"To: " + email + "\r\n" +
			"Subject: Email Verification\r\n\r\n" +
			messageBody,
	))

	_, err = service.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Verification email sent to %s", email)
	return nil
}

// GetTokenExpiry returns the token expiry time in minutes.
func (es *EmailService) GetTokenExpiry() int {
	return es.TokenExpiry
}

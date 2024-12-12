package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents the main configuration structure for the application.
type Config struct {
	AppPort        string `json:"app_port"`
	AppEnv         string `json:"app_env"`
	AppAPIBasePath string `json:"app_api_base_path"`

	JWTSecretKey string `json:"jwt_secret_key"`

	DBHost string `json:"db_host"`
	DBPort string `json:"db_port"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
	DBName string `json:"db_name"`

	GoogleOAuthAccessToken  string `json:"google_oauth_access_token"`
	GoogleOAuthRefreshToken string `json:"google_oauth_refresh_token"`
	GoogleClientID          string `json:"google_client_id"`
	GoogleClientSecret      string `json:"google_client_secret"`
	GoogleRedirectURI       string `json:"google_redirect_uri"`
	GoogleAuthURL           string `json:"google_auth_url"`
	GoogleTokenURL          string `json:"google_token_url"`

	EmailSender          string `json:"email_sender"`
	EmailTokenExpiry     int    `json:"email_token_expiry"`
	EmailVerificationURL string `json:"email_verification_url"`
}

// LoadConfig loads the configuration from the .env file or environment variables. It returns a pointer to the Config struct. If the .env file is not found, it logs a warning and uses environment variables instead. If any error occurs during the loading process, it logs an error and returns nil.
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Using environment variables.")
	}

	EmailTokenExpiry, err := strconv.Atoi(os.Getenv("EMAIL_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalf("Invalid value for SMT_EMAIL_TOKEN_EXPIRY: %v", err)
	}

	return &Config{
		AppPort:        os.Getenv("APP_PORT"),
		AppEnv:         os.Getenv("APP_ENV"),
		AppAPIBasePath: os.Getenv("APP_API_BASE_PATH"),

		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),

		DBHost: os.Getenv("POSTGRES_HOST"),
		DBPort: os.Getenv("POSTGRES_PORT"),
		DBUser: os.Getenv("POSTGRES_USER"),
		DBPass: os.Getenv("POSTGRES_PASSWORD"),
		DBName: os.Getenv("POSTGRES_DB"),

		GoogleOAuthAccessToken:  os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN"),
		GoogleOAuthRefreshToken: os.Getenv("GOOGLE_OAUTH_REFRESH_TOKEN"),
		GoogleClientID:          os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:      os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURI:       os.Getenv("GOOGLE_REDIRECT_URI"),
		GoogleAuthURL:           os.Getenv("GOOGLE_AUTH_URL"),
		GoogleTokenURL:          os.Getenv("GOOGLE_TOKEN_URL"),

		EmailTokenExpiry:     EmailTokenExpiry,
		EmailSender:          os.Getenv("EMAIL_SENDER"),
		EmailVerificationURL: os.Getenv("EMAIL_VERIFICATION_URL"),
	}
}

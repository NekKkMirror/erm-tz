package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/NekKkMirror/erm-tz/config"
	"github.com/NekKkMirror/erm-tz/internal/handler"
	"github.com/NekKkMirror/erm-tz/internal/repository"
	"github.com/NekKkMirror/erm-tz/internal/service"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App represents the main application structure.
// It holds the application configuration, database connection, and router.
type App struct {
	Config *config.Config
	DB     *sql.DB
	Router *mux.Router
}

// Init initializes a new App instance, setting up the configuration,
// database connection, and application routes.
func Init() *App {
	a := &App{}

	a.Config = config.LoadConfig()
	a.DB = a.initDB()
	a.Router = a.initRouter()

	return a
}

// initDB initializes and returns a connection to the PostgreSQL database.
// It uses the application's configuration to construct the connection string.
func (a *App) initDB() *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db", a.Config.DBPort, a.Config.DBUser, a.Config.DBPass, a.Config.DBName,
	))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
	return db
}

// initRouter initializes and sets up the HTTP router with the application's routes.
// It registers routes for user-related operations.
func (a *App) initRouter() *mux.Router {
	rout := mux.NewRouter()

	apiRouter := rout.PathPrefix(a.Config.AppAPIBasePath).Subrouter()

	userRepo := repository.NewUserRepository(a.DB)
	emailService := service.NewEmailService(a.Config)
	userService := service.NewUserService(userRepo, emailService)
	userHandler := handler.NewUserHandler(userService, emailService)

	handler.RegisterUsersRouter(apiRouter, userHandler)

	log.Println("Routes initialized successfully")
	return rout
}

// Run starts the HTTP server and logs the port the server is running on.
func (a *App) Run() {
	log.Println("Server running on port", a.Config.AppPort)
	log.Fatal(http.ListenAndServe(":"+a.Config.AppPort, a.Router))
}

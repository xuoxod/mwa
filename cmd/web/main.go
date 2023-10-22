package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/xuoxod/mwa/internal/config"
	"github.com/xuoxod/mwa/internal/handlers"
)

// Application configuration
var app config.AppConfig

// var templateData models.TemplateData
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	app.InProduction = false

	// App config middleware
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Session middleware
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// Config app level config with session
	app.Session = session

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	// Set the app level session
	app.Session = session

	go handlers.ListenToWsChannel()

	mux := routes()

	log.Println("Server running on port 8080")

	_ = http.ListenAndServe(":8080", mux)
}

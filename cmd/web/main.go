package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/xuoxod/mwa/internal/config"
	"github.com/xuoxod/mwa/internal/driver"
	"github.com/xuoxod/mwa/internal/envloader"
	"github.com/xuoxod/mwa/internal/handlers"
	"github.com/xuoxod/mwa/internal/helpers"
	"github.com/xuoxod/mwa/internal/models"
)

// Application configuration
var app config.AppConfig

// var templateData models.TemplateData
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	app.DBConnection = os.Getenv("DB_URL")

	err := envloader.LoadEnvVars()

	if err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	/* 	if err != nil {
		log.Fatal("Cannot connect to database! Dying ...")
	} */

	go handlers.ListenToWsChannel()

	mux := routes()

	log.Println("Server running on port 8080")

	_ = http.ListenAndServe(":8080", mux)
}

func run() (*driver.DB, error) {

	// Store a value in session
	gob.Register(models.Profile{})
	gob.Register(models.UserSettings{})
	gob.Register(models.Users{})
	gob.Register(models.User{})
	gob.Register(models.Signin{})
	gob.Register(models.Registration{})
	gob.Register(models.Authentication{})

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

	// Set the app level session
	app.Session = session

	var host string = os.Getenv("DB_HOST")
	var user string = os.Getenv("DB_USER")
	var password string = os.Getenv("DB_PASSWD")
	var dbname string = os.Getenv("DB_NAME")
	var port int = 5432

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := driver.ConnectSql(psqlInfo)

	if err != nil {
		log.Println("Cannot connect to database! Dying ...")
		return nil, err
	}

	fmt.Println("Connected to datastore")

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)
	helpers.NewHelpers(&app)

	return db, nil
}

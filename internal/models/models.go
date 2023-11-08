package models

import (
	"time"
)

// User registration data
type Registration struct {
	FirstName       string
	LastName        string
	Email           string
	Phone           string
	PasswordCreate  string
	PasswordConfirm string
}

type RegistrationErrData struct {
	Data map[string]string
}

// User signin data
type Signin struct {
	Email    string
	Password string
}

// All users
type Users struct {
	Collection map[string][]User
}

// User
type User struct {
	ID            string
	FirstName     string
	LastName      string
	Email         string
	Phone         string
	Password      string
	EmailVerified bool
	PhoneVerified bool
	AccessLevel   int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// User profile
type Profile struct {
	ID        int
	UserID    int
	UserName  string
	Image     byte
	Address   string
	City      string
	State     string
	Zipcode   string
	UpdatedAt time.Time
	CreatedAt time.Time
}

// User preferences
type Preferences struct {
	ID                       int
	UserID                   int
	EnablePublicProfile      bool
	EnableSmsNotifications   bool
	EnableEmailNotifications bool
	UpdatedAt                time.Time
	CreatedAt                time.Time
}

// Auth variable
type Authentication struct {
	Auth bool
}

package models

import "time"

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
	ID          string
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User profile
type Profile struct {
	ID          int
	UserID      int
	UserName    string
	DisplayName string
	ImageURL    string
	Address     string
	City        string
	State       string
	Zipcode     string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// User settings
type UserSettings struct {
	ID                int
	UserID            int
	ShowOnlineStatus  bool
	ShowPhone         bool
	ShowEmail         bool
	ShowNotifications bool
	ShowAddress       bool
	ShowCity          bool
	ShowState         bool
	ShowZipcode       bool
}

// Auth variable
type Authentication struct {
	Auth bool
}

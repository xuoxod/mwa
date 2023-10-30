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
	ID        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User profile
type Profile struct {
	UserID      int
	UserName    string
	DisplayName string
	ImageURL    string
	Address     string
	City        string
	State       string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// User settings
type UserSettings struct {
	UserID            int
	Online            bool
	ShowProfile       bool
	ShowUsername      bool
	ShowOnlineStatus  bool
	ShowAddress       bool
	ShowCity          bool
	ShowState         bool
	ShowDisplayName   bool
	ShowContactInfo   bool
	ShowPhone         bool
	ShowEmail         bool
	ShowNotifications bool
}

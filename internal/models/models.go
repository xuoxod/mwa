package models

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

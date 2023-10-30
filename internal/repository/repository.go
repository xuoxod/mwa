package repository

import "github.com/xuoxod/mwa/internal/models"

type DatabaseRepo interface {
	AllUsers() models.Users
	CreateUser(res models.Registration) (int, error)
	RemoveUser(id int) error
	GetUserByID(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateUser(models.User) (models.User, error)
	// UpdateUserProfile(userId int) (models.Profile, error)
	// UpdateUserSettings(u models.UserSettings) (models.UserSettings, error)
	Authenticate(email, testPassword string) (models.User, models.Profile, models.UserSettings, error)
}

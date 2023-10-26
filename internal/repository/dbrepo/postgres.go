package dbrepo

import "github.com/xuoxod/mwa/internal/models"

func (m *postgresDbRepo) AllUsers() models.Users {
	var users models.Users

	return users
}

func (m *postgresDbRepo) CreateUser(user models.Registration) (models.User, error) {
	var u models.User

	return u, nil
}

func (m *postgresDbRepo) RemoveUser(id int) error {

	return nil
}

func (m *postgresDbRepo) GetUserByID(id int) (models.User, error) {
	var user models.User

	return user, nil
}

func (m *postgresDbRepo) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	return user, nil
}

func (m *postgresDbRepo) UpdateUser(user models.User) (models.User, error) {
	var u models.User

	return u, nil
}

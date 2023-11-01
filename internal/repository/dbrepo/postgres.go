package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xuoxod/mwa/internal/helpers"
	"github.com/xuoxod/mwa/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDbRepo) AllUsers() models.Users {
	var users models.Users

	return users
}

func (m *postgresDbRepo) CreateUser(user models.Registration) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id, pId, sId int

	stmt := `insert into krxbyhhs.public.users(first_name, last_name, email, phone, password, created_at, updated_at) values($1,$2,$3,$4,$5,$6,$7) returning id`

	hashedPassword, hashPasswordErr := helpers.HashPassword(user.PasswordConfirm)

	if hashPasswordErr != nil {
		fmt.Println("Error hashing password: ", hashPasswordErr.Error())
		return 0, hashPasswordErr
	}

	row := m.DB.QueryRowContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		hashedPassword,
		time.Now(),
		time.Now(),
	)

	memberErr := row.Scan(&id)

	if memberErr != nil {
		fmt.Println("User Error: ", memberErr.Error())
		return 0, memberErr
	}

	stmt = `insert into krxbyhhs.public.profiles(user_id, created_at, updated_at, user_name, display_name, image_url, address, city, state, zipcode) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning user_id`

	row = m.DB.QueryRowContext(ctx, stmt, id, time.Now(), time.Now(), "user-name", "display-name", "image-url", "address", "city", "state", "zipcode")

	memberErr = row.Scan(&pId)

	if memberErr != nil {
		fmt.Println("Profile Error: ", memberErr.Error())
		return 0, memberErr
	}

	stmt = `insert into krxbyhhs.public.usersettings(user_id, created_at, updated_at, user_name) values($1,$2,$3) returning user_id`

	// Create unique username

	username := fmt.Sprintf("%s-%s", user.LastName, user.Email)

	row = m.DB.QueryRowContext(ctx, stmt, id, time.Now(), time.Now(), username)

	memberErr = row.Scan(&sId)

	if memberErr != nil {
		fmt.Println("Profile Error: ", memberErr.Error())
		return 0, memberErr
	}

	return id, nil
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

func (m *postgresDbRepo) Authenticate(email, testPassword string) (models.User, models.Profile, models.UserSettings, error) {
	var user models.User
	var profile models.Profile
	var userSettings models.UserSettings

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select u.id, u.first_name, u.last_name, u.email, u.phone, u.access_level, u.created_at, u.updated_at, u.password, p.user_name, p.display_name, p.image_url, p.address, p.city, p.state, p.zipcode, s.show_online_status, s.show_email, s.show_phone, s.show_notifications from users u inner join profiles p on p.user_id = u.id inner join usersettings s on s.user_id = u.id where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt, &user.Password, &profile.UserName, &profile.DisplayName, &profile.ImageURL, &profile.Address, &profile.City, &profile.State, &profile.Zipcode, &userSettings.ShowOnlineStatus, &userSettings.ShowEmail, &userSettings.ShowPhone, &userSettings.ShowNotifications)

	if err != nil {
		log.Printf("\n\tQuery error on table users\n\t%s\n", err.Error())
		return user, profile, userSettings, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("bcrypt error:\t", err.Error())

		return user, profile, userSettings, err
	} else if err != nil {
		log.Println("bcrypt error:\t", err.Error())

		return user, profile, userSettings, err
	}

	return user, profile, userSettings, nil
}

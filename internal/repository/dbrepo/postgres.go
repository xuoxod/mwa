package dbrepo

import (
	"context"
	"errors"
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
		fmt.Println("Error hashing password: ", errors.New(hashPasswordErr.Error()))
		return 0, errors.New(hashPasswordErr.Error())
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
		// fmt.Println("User Error: ", errors.New(memberErr.Error()))
		return 0, errors.New(memberErr.Error())
	}

	// Create unique username

	username := fmt.Sprintf("%s-%s", user.LastName, user.Email)

	stmt = `insert into krxbyhhs.public.profiles(user_id, created_at, updated_at, user_name, display_name, image_url, address, city, state, zipcode) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning user_id`

	row = m.DB.QueryRowContext(ctx, stmt, id, time.Now(), time.Now(), username, "display-name", "image-url", "address", "city", "state", "zipcode")

	memberErr = row.Scan(&pId)

	if memberErr != nil {
		// fmt.Println("Profile Error: ", errors.New(memberErr.Error()))
		return 0, errors.New(memberErr.Error())
	}

	stmt = `insert into krxbyhhs.public.usersettings(user_id, created_at, updated_at) values($1,$2,$3) returning user_id`

	row = m.DB.QueryRowContext(ctx, stmt, id, time.Now(), time.Now())

	memberErr = row.Scan(&sId)

	if memberErr != nil {
		// fmt.Println("Profile Error: ", errors.New(memberErr.Error()))
		return 0, errors.New(memberErr.Error())
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

// UpdateUser: update the user and profile
// @param models.User: first & last names, email, phone
// @param models.Profile: username, image url, address, city, state, zipcode
// @return models.User, models.Profile and error
func (m *postgresDbRepo) UpdateUser(user models.User, profile models.Profile) (models.User, models.Profile, error) {
	var u models.User
	var p models.Profile

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Update user table
	userQuery := `
		update users set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5 where email = $6 returning id, first_name, last_name, email, phone, updated_at
	`

	usersRows, usersRowsErr := m.DB.QueryContext(ctx, userQuery,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		time.Now(),
		user.Email,
	)

	if usersRowsErr != nil {
		return u, p, usersRowsErr
	}

	for usersRows.Next() {
		if err := usersRows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.UpdatedAt); err != nil {
			fmt.Printf("\tMember Row Scan Error: %s\n", err.Error())
			return u, p, err
		}
	}

	usersRerr := usersRows.Close()

	if usersRerr != nil {
		return u, p, usersRerr
	}

	if err := usersRows.Err(); err != nil {
		return u, p, err
	}

	// Update profile table

	profilesQuery := `
	update profiles set user_name = $1, image_url = $2, address = $3, city = $4, state = $5, zipcode = $6, updated_at = $7 where user_id = $8 returning user_name, image_url, address, city, state, zipcode, updated_at`

	profileRows, profileErr := m.DB.QueryContext(ctx, profilesQuery,
		profile.UserName,
		profile.ImageURL,
		profile.Address,
		profile.City,
		profile.State,
		profile.Zipcode,
		time.Now(),
		u.ID,
	)

	if profileErr != nil {
		return u, p, profileErr
	}

	for profileRows.Next() {
		if err := profileRows.Scan(&p.UserName, &p.ImageURL, &p.Address, &p.City, &p.State, &p.Zipcode, &p.UpdatedAt); err != nil {
			return u, p, err
		}
	}

	profileRerr := profileRows.Close()

	if profileRerr != nil {
		return u, p, profileRerr
	}

	if err := profileRows.Err(); err != nil {
		return u, p, err
	}

	return u, p, nil
}

func (m *postgresDbRepo) Authenticate(email, testPassword string) (models.User, models.Profile, models.UserSettings, error) {
	var user models.User
	var profile models.Profile
	var userSettings models.UserSettings

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select u.id, u.first_name, u.last_name, u.email, u.phone, u.access_level, u.created_at, u.updated_at, u.password, p.user_name, p.display_name, p.image_url, p.address, p.city, p.state, p.zipcode, s.show_online_status, s.show_email, s.show_phone, s.enable_sms_notifications, s.enable_email_notifications from users u inner join profiles p on p.user_id = u.id inner join usersettings s on s.user_id = u.id where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt, &user.Password, &profile.UserName, &profile.DisplayName, &profile.ImageURL, &profile.Address, &profile.City, &profile.State, &profile.Zipcode, &userSettings.ShowOnlineStatus, &userSettings.ShowEmail, &userSettings.ShowPhone, &userSettings.EnableSmsNotifications, &userSettings.EnableEmailNotifications)

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

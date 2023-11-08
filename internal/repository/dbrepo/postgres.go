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

	var userId, profileId, prefId int

	// Insert user
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

	rowErr := row.Scan(&userId)

	if rowErr != nil {
		fmt.Println("User Error: ", rowErr.Error())
		return 0, rowErr
	}

	// Insert profile

	// Create unique username
	username := fmt.Sprintf("%s-%s", user.LastName, user.Email)

	stmt = `insert into krxbyhhs.public.profiles(user_id, created_at, updated_at, user_name, address, city, state, zipcode) values($1,$2,$3,$4,$5,$6,$7,$8) returning id`

	row = m.DB.QueryRowContext(ctx, stmt, userId, time.Now(), time.Now(), username, "Enter address", "Enter city", "Enter state", "Enter zipcode")

	rowErr = row.Scan(&profileId)

	if rowErr != nil {
		fmt.Println("Profile Error: ", rowErr.Error())
		return 0, rowErr
	}

	// Insert preferences
	stmt = `insert into krxbyhhs.public.preferences(user_id, created_at, updated_at) values($1,$2,$3) returning id`

	row = m.DB.QueryRowContext(ctx, stmt, userId, time.Now(), time.Now())

	rowErr = row.Scan(&prefId)

	if rowErr != nil {
		fmt.Println("Preferences Error: ", rowErr.Error())
		return 0, rowErr
	}

	return userId, nil
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

// UpdateSettings: update the user settings
// @param models.User: first & last names, email, phone
// @param models.Profile: username, image url, address, city, state, zipcode
// @return models.User, models.Profile and error
func (m *postgresDbRepo) UpdatePreferences(preferences models.Preferences) (models.Preferences, error) {
	var p models.Preferences

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Update user table
	preferencesQuery := `
		update preferences set enable_sms_nots = $1, enable_email_nots = $2, enable_public_profile = $3, updated_at = $4 where user_id = $5 returning id, user_id, enable_sms_nots, enable_email_nots, enable_public_profile, updated_at
	`

	preferencesRows, preferencesRowsErr := m.DB.QueryContext(ctx, preferencesQuery,
		preferences.EnableSmsNotifications,
		preferences.EnableEmailNotifications,
		preferences.EnablePublicProfile,
		time.Now(),
		preferences.UserID,
	)

	if preferencesRowsErr != nil {
		log.Println("preferences row err: ", preferencesRowsErr)
		return p, preferencesRowsErr
	}

	for preferencesRows.Next() {
		if err := preferencesRows.Scan(&p.ID, &p.UserID, &p.EnableSmsNotifications, &p.EnableEmailNotifications, &p.EnablePublicProfile, &p.UpdatedAt); err != nil {
			log.Println("preferences rows scan err ", err)
			return p, err
		}
	}

	preferencesRerr := preferencesRows.Close()

	if preferencesRerr != nil {
		log.Println("preferences rerr err: ", preferencesRerr)
		return p, preferencesRerr
	}

	if err := preferencesRows.Err(); err != nil {
		log.Println("preferences row err(): ", err)
		return p, err
	}

	log.Println("Returning updated user preferences")
	log.Printf("\tPreferences:\n\t%v\n\n", p)

	return p, nil
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
	update profiles set user_name = $1, address = $2, city = $3, state = $4, zipcode = $5, updated_at = $6 where user_id = $7 returning user_name, address, city, state, zipcode, updated_at`

	profileRows, profileErr := m.DB.QueryContext(ctx, profilesQuery,
		profile.UserName,
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
		if err := profileRows.Scan(&p.UserName, &p.Address, &p.City, &p.State, &p.Zipcode, &p.UpdatedAt); err != nil {
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

func (m *postgresDbRepo) Authenticate(email, testPassword string) (models.User, models.Profile, models.Preferences, error) {
	var user models.User
	var profile models.Profile
	var preferences models.Preferences

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select u.id, u.first_name, u.last_name, u.email, u.phone, u.access_level, u.created_at, u.updated_at, u.password, u.email_verified, u.phone_verified, p.user_name, p.address, p.city, p.state, p.zipcode, s.id, s.user_id, s.enable_sms_nots, s.enable_email_nots, s.enable_public_profile from users u inner join profiles p on p.user_id = u.id inner join preferences s on s.user_id = u.id where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt, &user.Password, &user.EmailVerified, &user.PhoneVerified, &profile.UserName, &profile.Address, &profile.City, &profile.State, &profile.Zipcode, &preferences.ID, &preferences.UserID, &preferences.EnableSmsNotifications, &preferences.EnableEmailNotifications, &preferences.EnablePublicProfile)

	if err != nil {
		log.Printf("\n\tQuery error on table users\n\t%s\n", err.Error())
		return user, profile, preferences, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("bcrypt mismatched error:\t", err.Error())

		return user, profile, preferences, err
	} else if err != nil {
		log.Println("bcrypt error:\t", err.Error())

		return user, profile, preferences, err
	}

	return user, profile, preferences, nil
}

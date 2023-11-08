package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/xuoxod/mwa/internal/forms"
	"github.com/xuoxod/mwa/internal/helpers"
	"github.com/xuoxod/mwa/internal/models"
	"github.com/xuoxod/mwa/internal/render"
)

// @desc        User dashboard
// @route       GET /user
// @access      Private
func (m *Respository) UserDashboard(w http.ResponseWriter, r *http.Request) {
	auth, authOk := m.App.Session.Get(r.Context(), "auth").(models.Authentication)
	profile, profileOk := m.App.Session.Get(r.Context(), "profile").(models.Profile)
	preferences, preferencesOk := m.App.Session.Get(r.Context(), "preferences").(models.Preferences)
	user, userOk := m.App.Session.Get(r.Context(), "user_id").(models.User)

	if !authOk {
		log.Println("Cannot get auth data from session")
		m.App.ErrorLog.Println("Can't get auth data from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get auth data from session")
		http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
		return
	}

	if !profileOk {
		log.Println("Cannot get profile data from session")
		m.App.ErrorLog.Println("Can't get profile data from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get profile data from session")
		http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
		return
	}

	if !preferencesOk {
		log.Println("Cannot get preferences data from session")
		m.App.ErrorLog.Println("Can't get preferences data from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get preferences data from session")
		http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
		return
	}

	if !userOk {
		log.Println("Cannot get user_id data from session")
		m.App.ErrorLog.Println("Can't get user_id data from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get user_id data from session")
		http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]string)
	data["dashboard"] = fmt.Sprintf("%t", true)

	obj := make(map[string]interface{})
	obj["auth"] = auth
	obj["profile"] = profile
	obj["preferences"] = preferences
	obj["user"] = user
	obj["csrftoken"] = nosurf.Token(r)
	obj["title"] = "Dashboard"
	obj["form"] = forms.New(nil)

	err := render.RenderPageWithContext(w, "user/dashboard.jet", data, obj)

	if err != nil {
		log.Println(err.Error())
	}
}

// @desc        Update profile
// @route       POST /user/profile
// @access      Private
func (m *Respository) ProfilePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post profile")

	err := r.ParseForm()

	if err != nil {
		fmt.Printf("\n\tError parsing user profile form")
		helpers.ServerError(w, err)
		return
	}

	parsedProfile := models.Profile{
		UserName: r.Form.Get("uname"),
		Address:  r.Form.Get("address"),
		City:     r.Form.Get("city"),
		State:    r.Form.Get("state"),
		Zipcode:  r.Form.Get("zipcode"),
	}

	parsedUser := models.User{
		FirstName: r.Form.Get("fname"),
		LastName:  r.Form.Get("lname"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	// form validation
	form := forms.New(r.PostForm)
	form.IsEmail("email")
	form.IsUrl("iurl")
	form.Required("fname", "lname", "email", "phone")

	obj := make(map[string]interface{})

	if !form.Valid() {
		fmt.Println("Form Errors: ", form.Errors)
		obj["profileform"] = parsedProfile
		obj["userform"] = parsedUser
		obj["ok"] = false

		if form.Errors.Get("email") != "" {
			obj["email"] = form.Errors.Get("email")
		}

		if form.Errors.Get("iurl") != "" {
			obj["iurl"] = form.Errors.Get("iurl")
		}
		if form.Errors.Get("fname") != "" {
			obj["fname"] = form.Errors.Get("fname")
		}
		if form.Errors.Get("lname") != "" {
			obj["lname"] = form.Errors.Get("lname")
		}
		if form.Errors.Get("email") != "" {
			obj["email"] = form.Errors.Get("email")
		}
		if form.Errors.Get("phone") != "" {
			obj["phone"] = form.Errors.Get("phone")
		}

		out, err := json.MarshalIndent(obj, "", " ")

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		num, rErr := w.Write(out)

		if rErr != nil {
			log.Println(err)
		}

		log.Printf("Response Writer's returned integer: %d\n", num)
	} else {
		// Update user and their profile then return it
		updatedUser, updatedProfile, err := m.DB.UpdateUser(parsedUser, parsedProfile)

		if err != nil {
			fmt.Println(err)
		}

		// replace user_id and profile in the session manager
		m.App.Session.Remove(r.Context(), "user_id")
		m.App.Session.Remove(r.Context(), "profile")

		m.App.Session.Put(r.Context(), "user_id", updatedUser)
		m.App.Session.Put(r.Context(), "profile", updatedProfile)

		obj["ok"] = true

		out, err := json.MarshalIndent(obj, "", " ")

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		num, rErr := w.Write(out)

		if rErr != nil {
			log.Println(err)
		}

		log.Printf("Response Writer's returned integer: %d\n", num)
	}
}

// @desc        Update settings
// @route       POST /user/settings
// @access      Private
func (m *Respository) PreferencesPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post settings")
	obj := make(map[string]interface{})
	preferences, preferencesOk := m.App.Session.Get(r.Context(), "preferences").(models.Preferences)

	if !preferencesOk {
		log.Println("Cannot get preferences data from session")
		m.App.ErrorLog.Println("Can't get preferences data from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get preferences data from session")
		http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()

	if err != nil {
		fmt.Printf("\n\tError parsing user preferences form")
		helpers.ServerError(w, err)
		return
	}

	var parsedPreferences models.Preferences

	parsedPreferences.ID = preferences.ID
	parsedPreferences.UserID = preferences.UserID

	for key := range r.Form {
		if key == "enable-public-profile" {
			parsedPreferences.EnablePublicProfile = true
		}

		if key == "enable-sms-notifications" {
			parsedPreferences.EnableSmsNotifications = true
		}

		if key == "enable-email-notifications" {
			parsedPreferences.EnableEmailNotifications = true
		}

	}

	log.Printf("\n\tParsed Settings Form: \n\t%v\n\n", parsedPreferences)

	// Update user and their profile then return it
	updatedPreferences, err := m.DB.UpdatePreferences(parsedPreferences)

	if err != nil {
		fmt.Println(err)

		obj["ok"] = false

		out, err := json.MarshalIndent(obj, "", " ")

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, rErr := w.Write(out)

		if rErr != nil {
			log.Println(err)
		}
		return
	}

	m.App.Session.Remove(r.Context(), "preferences")
	m.App.Session.Put(r.Context(), "preferences", updatedPreferences)

	obj["ok"] = true

	out, err := json.MarshalIndent(obj, "", " ")

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, rErr := w.Write(out)

	if rErr != nil {
		log.Println(err)
	}
}

// @desc        Signout user
// @route       GET /user/signout
// @access      Private
func (m *Respository) SignOut(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

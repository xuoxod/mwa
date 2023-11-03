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
	settings, settingsOk := m.App.Session.Get(r.Context(), "settings").(models.UserSettings)
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

	if !settingsOk {
		log.Println("Cannot get settings data from session")
		m.App.ErrorLog.Println("Can't get settings data from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get settings data from session")
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
	obj["settings"] = settings
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
		ImageURL: r.Form.Get("iurl"),
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

		obj["user"] = parsedUser
		obj["profile"] = parsedProfile
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

// @desc        Signout user
// @route       GET /user/signout
// @access      Private
func (m *Respository) SignOut(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

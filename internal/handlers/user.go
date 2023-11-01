package handlers

import (
	"fmt"
	"log"
	"net/http"

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
	obj["title"] = "Dashboard"

	err := render.RenderPageWithContext(w, "user/dashboard.jet", data, obj)

	if err != nil {
		log.Println(err.Error())
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

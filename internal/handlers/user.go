package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xuoxod/mwa/internal/models"
	"github.com/xuoxod/mwa/internal/render"
)

func (m *Respository) UserDashboard(w http.ResponseWriter, r *http.Request) {
	auth, ok := m.App.Session.Get(r.Context(), "auth").(models.Authentication)

	if !ok {
		log.Println("Cannot get auth data from session")
		m.App.ErrorLog.Println("Can't get auth data from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get auth data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]string)
	data["title"] = "Dashboard"
	data["dashboard"] = fmt.Sprintf("%t", true)

	obj := make(map[string]interface{})
	obj["auth"] = auth

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

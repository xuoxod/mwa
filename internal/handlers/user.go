package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

func (m *Respository) UserDashboard(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)
	vars.Set("title", "Dashboard")
	vars.Set("dashboard", true)

	err := RenderPage(w, "user/dashboard.jet", vars)

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

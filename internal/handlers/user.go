package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

func (m *Respository) Dashboard(w http.ResponseWriter, r *http.Request) {
	// remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	vars := make(jet.VarMap)
	vars.Set("title", "Dashboard")
	// vars.Set("ipaddress", remoteIp)
	vars.Set("dashboard", true)

	err := RenderPage(w, "user/dashboard.jet", vars)

	if err != nil {
		log.Println(err.Error())
	}
}

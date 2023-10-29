package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

func (m *Respository) Dashboard(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)
	vars.Set("title", "Public")
	vars.Set("public", true)

	err := RenderPage(w, "public/dashboard.jet", vars)

	if err != nil {
		log.Println(err.Error())
	}
}

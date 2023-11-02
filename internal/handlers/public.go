package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xuoxod/mwa/internal/render"
)

func (m *Respository) Dashboard(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["title"] = "Public"
	data["public"] = fmt.Sprintf("%t", true)

	err := render.RenderPage(w, "public/dashboard.jet", data)

	if err != nil {
		log.Println(err.Error())
	}
}

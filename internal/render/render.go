package render

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/CloudyKit/jet"
)

// var views = jet.NewSet(
// jet.NewOSFileSystemLoader("./views"),
// jet.InDevelopmentMode(),
// )

var root, _ = os.Getwd()
var views = jet.NewHTMLSet(filepath.Join(root, "views"))

func RenderPage(w http.ResponseWriter, tmpl string, data map[string]string) error {
	view, err := views.GetTemplate(tmpl)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	vmap := make(jet.VarMap)

	for k, v := range data {
		fmt.Println(k, v)
		vmap.Set(k, v)
	}

	err = view.Execute(w, vmap, nil)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func RenderPageWithContext(w http.ResponseWriter, tmpl string, data map[string]string, obj map[string]interface{}) error {
	view, err := views.GetTemplate(tmpl)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	vmap := make(jet.VarMap)

	for k, v := range data {
		vmap.Set(k, v)
	}

	err = view.Execute(w, vmap, obj)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

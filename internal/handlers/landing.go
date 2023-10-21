package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/justinas/nosurf"
	"github.com/xuoxod/mwa/internal/config"
	"github.com/xuoxod/mwa/internal/forms"
	"github.com/xuoxod/mwa/internal/helpers"
	"github.com/xuoxod/mwa/internal/models"
	"github.com/xuoxod/mwa/pkg/utils"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.InDevelopmentMode(),
)

// Repo the repository used by the handlers
var Repo *Respository

// Repository the Repository type
type Respository struct {
	App *config.AppConfig
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig) *Respository {
	return &Respository{
		App: a,
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Respository) {
	Repo = r
}

func (m *Respository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	vars := make(jet.VarMap)
	vars.Set("title", "Home")
	vars.Set("headingOne", `Welcome to Awesome Web App`)
	vars.Set("statement", `We don't Fuck around ... Either sign in to start using the site or register first then sign in.`)

	var emptySigninForm models.Signin

	data := make(map[string]interface{})
	data["csrftoken"] = nosurf.Token(r)
	data["signinform"] = emptySigninForm
	data["form"] = forms.New(nil)

	err := RenderPageWithContext(w, "landing/home.jet", vars, data)

	if err != nil {
		log.Println(err.Error())
	}
}

func (m *Respository) About(w http.ResponseWriter, r *http.Request) {
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	vars := make(jet.VarMap)
	vars.Set("title", "About")
	vars.Set("appname", `Awesome Web App`)
	vars.Set("appver", `1.0`)
	vars.Set("appdate", utils.DateTimeStamp())
	vars.Set("remoteip", remoteIp)

	err := RenderPage(w, "landing/about.jet", vars)

	if err != nil {
		log.Println(err.Error())
	}
}

func (m *Respository) Register(w http.ResponseWriter, r *http.Request) {
	// remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	var emptyRegistrationForm models.Registration
	vars := make(jet.VarMap)
	vars.Set("title", "Registration")

	data := make(map[string]interface{})
	data["csrftoken"] = nosurf.Token(r)
	data["registrationform"] = emptyRegistrationForm
	data["form"] = forms.New(nil)

	err := RenderPageWithContext(w, "landing/register.jet", vars, data)

	if err != nil {
		log.Println(err.Error())
	}
}

func (m *Respository) PostRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post Register")
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	registration := models.Registration{
		FirstName:       r.Form.Get("fname"),
		LastName:        r.Form.Get("lname"),
		Email:           r.Form.Get("email"),
		Phone:           r.Form.Get("phone"),
		PasswordCreate:  r.Form.Get("pwd1"),
		PasswordConfirm: r.Form.Get("pwd2"),
	}

	fmt.Println("Registration posted")

	form := forms.New(r.PostForm)
	form.Required("fname", "lname", "email", "phone", "pwd1", "pwd2")
	form.MinLength("fname", 2, r)
	form.MinLength("lname", 2, r)
	form.IsEmail("email")
	form.PasswordsMatch("pwd1", "pwd2", r)

	if !form.Valid() {
		fmt.Println(form.Errors)

		vars := make(jet.VarMap)
		vars.Set("title", "Registration")

		data := make(map[string]interface{})
		data["csrftoken"] = nosurf.Token(r)
		data["registrationform"] = registration
		data["form"] = form

		err := RenderPageWithContext(w, "landing/register.jet", vars, data)

		if err != nil {
			log.Println(err.Error())
		}
	} else {
		// Create new user in the database

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (m *Respository) PostSignin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post Signin")
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	signinform := models.Signin{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	fmt.Println("Signin posted")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		fmt.Println(form.Errors)

		vars := make(jet.VarMap)
		vars.Set("title", "Home")
		vars.Set("headingOne", `Welcome to Awesome Web App`)
		vars.Set("statement", `We don't Fuck around ... Either sign in to start using the site or register first then sign in.`)

		data := make(map[string]interface{})
		data["csrftoken"] = nosurf.Token(r)
		data["signinform"] = signinform
		data["form"] = form

		err := RenderPageWithContext(w, "landing/home.jet", vars, data)

		if err != nil {
			log.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
	}
}

func RenderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = view.Execute(w, data, nil)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func RenderPageWithContext(w http.ResponseWriter, tmpl string, data jet.VarMap, obj map[string]interface{}) error {
	view, err := views.GetTemplate(tmpl)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = view.Execute(w, data, &obj)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/justinas/nosurf"
	"github.com/xuoxod/mwa/internal/config"
	"github.com/xuoxod/mwa/internal/driver"
	"github.com/xuoxod/mwa/internal/forms"
	"github.com/xuoxod/mwa/internal/helpers"
	"github.com/xuoxod/mwa/internal/models"
	"github.com/xuoxod/mwa/internal/render"
	"github.com/xuoxod/mwa/internal/repository"
	"github.com/xuoxod/mwa/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Respository

// Repository the Repository type
type Respository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Respository {
	return &Respository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Respository) {
	Repo = r
	render.InitViews()
}

func (m *Respository) Home(w http.ResponseWriter, r *http.Request) {
	// vars := make(jet.VarMap)
	// vars := Vmap
	// vars.Set("title", "Home")
	var emptySigninForm models.Signin
	data := make(map[string]string)

	obj := make(map[string]interface{})
	obj["csrftoken"] = nosurf.Token(r)
	obj["signinform"] = emptySigninForm
	obj["form"] = forms.New(nil)
	obj["title"] = "Home"

	err := render.RenderPageWithContext(w, "landing/home.jet", data, obj)

	if err != nil {
		log.Println(err.Error())
	}
}

func (m *Respository) About(w http.ResponseWriter, r *http.Request) {
	// vars := make(jet.VarMap)
	// vars.Set("title", "About")
	// vars.Set("appname", `Awesome Web App`)
	// vars.Set("appver", `1.0`)
	// vars.Set("appdate", utils.DateTimeStamp())

	data := make(map[string]string)
	data["title"] = "About"

	err := render.RenderPage(w, "landing/about.jet", data)

	if err != nil {
		log.Println(err.Error())
	}
}

func (m *Respository) Register(w http.ResponseWriter, r *http.Request) {
	regData, regDataOk := m.App.Session.Get(r.Context(), "reg-error").(models.RegistrationErrData)

	if !regDataOk {
		log.Println("Cannot get reg-error data from session")
		// m.App.ErrorLog.Println("Can't get reg-error data from the session")
		// m.App.Session.Put(r.Context(), "error", "Can't get reg-error data from session")
		// http.Redirect(w, r, "/register", http.StatusTemporaryRedirect)
		// return
	}

	var emptyRegistrationForm models.Registration
	// data := make(jet.VarMap)
	// data.Set("title", "Registration")
	data := make(map[string]string)

	obj := make(map[string]interface{})
	obj["csrftoken"] = nosurf.Token(r)
	obj["registrationform"] = emptyRegistrationForm
	obj["form"] = forms.New(nil)
	obj["title"] = "Registration"

	if regDataOk {
		data["error"] = regData.Data["error"]
		obj["type"] = regData.Data["type"]
		obj["msg"] = regData.Data["msg"]
	}

	err := render.RenderPageWithContext(w, "landing/register.jet", data, obj)

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

		// vars := make(jet.VarMap)
		// vars.Set("title", "Registration")

		data := make(map[string]string)
		data["title"] = "Registration'"

		obj := make(map[string]interface{})
		obj["csrftoken"] = nosurf.Token(r)
		obj["registrationform"] = registration
		obj["form"] = form

		err := render.RenderPageWithContext(w, "landing/register.jet", data, obj)

		if err != nil {
			log.Println(err.Error())
		}
	} else {
		// Create new user in the database
		// ERROR: duplicate key value violates unique constraint "users_un" (SQLSTATE 23505)
		userId, err := m.DB.CreateUser(registration)

		if err != nil {
			sErr := err.Error()
			uniqueErr := strings.HasSuffix(sErr, "(SQLSTATE 23505)")

			if uniqueErr {
				fmt.Println("Record already exists")
				var registrationErrData models.RegistrationErrData

				regErrData := make(map[string]string)
				regErrData["title"] = "Home"
				regErrData["error"] = "Authentication Error"
				regErrData["type"] = "error"
				regErrData["msg"] = "Account already exists"

				registrationErrData.Data = regErrData

				m.App.Session.Put(r.Context(), "reg-error", registrationErrData)

			}

			m.App.Session.Put(r.Context(), "error", "Error registering user")
			vars := make(jet.VarMap)
			vars.Set("title", "Registration")

			data := make(map[string]interface{})
			data["csrftoken"] = nosurf.Token(r)
			data["registrationform"] = registration
			data["form"] = form

			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		if userId > 0 {
			fmt.Println("User created successfully")
		}

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

		// vars := make(jet.VarMap)
		// vars.Set("title", "Home")

		data := make(map[string]string)
		data["title"] = "Home"

		obj := make(map[string]interface{})
		obj["csrftoken"] = nosurf.Token(r)
		obj["signinform"] = signinform
		obj["form"] = form

		err := render.RenderPageWithContext(w, "landing/home.jet", data, obj)

		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	// Authenticate user

	u, p, s, err := m.DB.Authenticate(signinform.Email, signinform.Password)

	if err != nil {
		fmt.Println("Authentication Error:\t", err.Error())
		// vars := make(jet.VarMap)
		// vars.Set("title", "Home")
		// vars.Set("error", "Authentication Error")

		vars := make(map[string]string)
		vars["title"] = "Home"
		vars["error"] = "Authentication Error"

		data := make(map[string]interface{})
		data["type"] = "error"
		data["msg"] = "Invalid Signin Credentials"
		data["signinform"] = signinform
		data["form"] = form
		data["csrftoken"] = nosurf.Token(r)

		err := render.RenderPageWithContext(w, "landing/home.jet", vars, data)

		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	var user models.User
	var profile models.Profile
	var settings models.UserSettings

	user = u
	profile = p
	settings = s

	// Put user in session
	m.App.Session.Put(r.Context(), "user_id", user)
	m.App.Session.Put(r.Context(), "profile", profile)
	m.App.Session.Put(r.Context(), "settings", settings)

	if user.AccessLevel == 1 {
		m.App.Session.Put(r.Context(), "admin_id", user)
	}

	var auth models.Authentication
	auth.Auth = true
	m.App.Session.Put(r.Context(), "auth", auth)
	http.Redirect(w, r, "/user", http.StatusSeeOther)

}

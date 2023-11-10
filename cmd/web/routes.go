package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xuoxod/mwa/internal/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Compress(5))
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)
	// mux.Use(RecoverPanic)
	mux.Use(WriteToConsole)
	mux.Use(middleware.NoCache)
	mux.Use(NoSurf)

	mux.Route("/", func(mux chi.Router) {
		mux.Use(Unauth)
		mux.Get("/", handlers.Repo.Home)
		mux.Get("/about", handlers.Repo.About)
		mux.Get("/register", handlers.Repo.Register)
		mux.Post("/register", handlers.Repo.PostRegister)
		mux.Post("/", handlers.Repo.PostSignin)
	})

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.UserDashboard)
		mux.Get("/signout", handlers.Repo.SignOut)
		mux.Post("/profile", handlers.Repo.ProfilePost)
		mux.Post("/settings", handlers.Repo.PreferencesPost)
		mux.Get("/settings", handlers.Repo.Settings)
		mux.Get("/email/verify", handlers.Repo.VerifyEmail)
		mux.Post("/email/verify", handlers.Repo.VerifyEmailPost)
		mux.Get("/phone/verify", handlers.Repo.VerifyPhone)
		mux.Post("/phone/verify", handlers.Repo.VerifyPhonePost)
	})

	/* 	mux.Route("/chat", func(mux chi.Router) {
	   		mux.Use(Auth)
	   		mux.Get("/", handlers.Repo.WsEndpoint)
	   	})
	*/
	fileserver := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileserver))

	// router.Handle("/static/*", http.StripPrefix("/static/", fs))

	return mux
}

package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xuoxod/mwa/internal/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	// mux.Use(middleware.Compress(5))
	mux.Use(middleware.Recoverer)
	// mux.Use(RecoverPanic)
	mux.Use(WriteToConsole)
	mux.Use(middleware.NoCache)
	mux.Use(NoSurf)
	// mux.Use(session.LoadAndSave)

	fileserver := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	mux.Route("/", func(mux chi.Router) {
		mux.Get("/", handlers.Repo.Home)
		mux.Get("/about", handlers.Repo.About)
		mux.Get("/register", handlers.Repo.Register)
		mux.Post("/register", handlers.Repo.PostRegister)
		mux.Post("/", handlers.Repo.PostSignin)
	})

	mux.Route("/user", func(mux chi.Router) {
		mux.Get("/", handlers.Repo.Dashboard)
	})

	mux.Route("/ws", func(mux chi.Router) {
		// mux.Use(HijackThis)
		mux.Get("/", handlers.Repo.WsEndpoint)
	})

	return mux
}

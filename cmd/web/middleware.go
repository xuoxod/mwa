package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/xuoxod/mwa/internal/helpers"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.RemoteAddr
		host := r.Host
		path := r.URL.Path
		method := r.Method
		protocol := r.Proto
		protocolMajor := r.ProtoMajor
		protocolMinor := r.ProtoMinor

		fmt.Printf("Middleware: type is %T\n", w)

		if path == "/ws" {
			w.Header().Set("Content", "text/plain")
		}

		fmt.Printf("\nPage Hit\n\tHost: %v\n\taddress: %v\n\tPath: %v\n\tMethod: %v\n\tProtocol: %v\n\t\tMajor: %v\n\t\tMinor: %v\n", host, remoteAddr, path, method, protocol, protocolMajor, protocolMinor)
		next.ServeHTTP(w, r)
	})
}

func HijackThis(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// w.Header().Set("Connection", "Upgrade")
		// w.Header().Set("Accept-Encoding", "gzip, deflate, br")
		// w.Header().Set("Connection", "keep-alive")
		// w.Header().Set("Cache-Control", "max-age=0")

		// Hijacking the connection
		h, _ := w.(http.Hijacker)

		conn, br, err := h.Hijack()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		responseBody := "Http response from hijacked connection"
		hr := http.Response{
			Status:           "200 OK",
			StatusCode:       200,
			Proto:            "HTTP/1.1",
			ProtoMajor:       1,
			ProtoMinor:       1,
			Header:           make(http.Header, 0),
			Body:             ioutil.NopCloser(bytes.NewBufferString(responseBody)),
			ContentLength:    int64(len(responseBody)),
			TransferEncoding: nil,
			Close:            false,
			Uncompressed:     false,
			Trailer:          nil,
			Request:          r,
			TLS:              nil,
		}

		// Writing response
		err = hr.Write(br)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Sending EOF to allow io.ReadAll(resp.Body) without blocking
		if v, ok := conn.(interface{ CloseWrite() error }); ok {
			err = v.CloseWrite()
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RecoverPanic recovers from a panic
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Check if there has been a panic
			if err := recover(); err != nil {
				// return a 500 Internal Server response
				helpers.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Sign in first")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Unauth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "warning", "Resource not found")
			http.Redirect(w, r, "/user", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAdmin(r) {
			session.Put(r.Context(), "warning", "Access Restricted")
			http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

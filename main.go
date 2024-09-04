package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			return
		}
	})
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		_, err := fmt.Fprint(os.Stderr, err)
		if err != nil {
			return
		}
		return
	}
}

// todo:

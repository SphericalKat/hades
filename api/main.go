package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"hades-2.0/api/middleware"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.Use(middleware.JwtAuthentication)

	n := negroni.Classic()
	n.UseHandler(r)

	err := http.ListenAndServe(":4000", n)
	if err != nil {
		log.Panic(err)
	}
}

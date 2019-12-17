package main

import (
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/hades-2/pkg/participant"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.Use(middleware.JwtAuthentication)

	n := negroni.Classic()
	n.UseHandler(r)

	db := connectDb()

	log.Println("Listening on port 4000")

	err := http.ListenAndServe(":4000", n)
	if err != nil {
		log.Panic(err)
	}
}

func connectDb() *gorm.DB {
	conn, err := pq.ParseURL(os.Getenv("DB_URI"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("DEBUG") == "true" {
		db = db.Debug()
	}
	db.AutoMigrate(&participant.Participant{}, &event.Event{})

	return db
}

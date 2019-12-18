package main

import (
	"fmt"
	"github.com/ATechnoHazard/hades-2/api/handler"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/hades-2/pkg/participant"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	negronilogrus "github.com/meatballhat/negroni-logrus"
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

func initNegroni() *negroni.Negroni {
	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewStatic(http.Dir("public")))
	return n
}

func main() {
	r := mux.NewRouter()

	n := initNegroni()
	n.UseHandler(r)

	db := connectDb()

	partRepo := participant.NewPostgresRepo(db)
	//eventRepo := event.NewPostgresRepo(db)

	partSvc := participant.NewParticipantService(partRepo)
	//eventSvc := event.NewEventService(eventRepo)
	handler.MakeParticipantHandler(r, partSvc)

	//_ = eventSvc.SaveEvent(&event.Event{
	//	ID:                    1,
	//	ClubName:              "DSC",
	//	Name:                  "SOME",
	//	Budget:                "69",
	//	Description:           "OWO",
	//	Category:              "NO",
	//	Venue:                 "YES",
	//	Attendance:            "DF",
	//	ExpectedParticipants:  "sDF",
	//	PROrequest:            "sdfsdf",
	//	CampusEngineerRequest: "sdfsdf",
	//	Duration:              "SFD",
	//	Status:                "sdfsdf",
	//	ToDate:                time.Now(),
	//	FromDate:              time.Now(),
	//	ToTime:                time.Now().Add(time.Hour),
	//	FromTime:              time.Now(),
	//})


	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Println("Listening on port " + port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), n)
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

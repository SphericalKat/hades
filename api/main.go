package main

import (
	"fmt"
	"github.com/ATechnoHazard/hades-2/api/handler"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/hades-2/pkg/organization"
	"github.com/ATechnoHazard/hades-2/pkg/participant"
	"github.com/ATechnoHazard/hades-2/pkg/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
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

	db.AutoMigrate(&entities.Participant{}, &entities.Event{}, &entities.Organization{}, &entities.User{})
	return db
}

func main() {
	r := httprouter.New()

	n := initNegroni()
	n.UseHandler(r)
	db := connectDb()

	partRepo := participant.NewPostgresRepo(db)
	eventRepo := event.NewPostgresRepo(db)
	orgRepo := organization.NewPostgresRepo(db)
	userRepo := user.NewPostgresRepo(db)

	partSvc := participant.NewParticipantService(partRepo)
	eventSvc := event.NewEventService(eventRepo)
	orgSvc := organization.NewOrganizationService(orgRepo)
	userSvc := user.NewUserService(userRepo)

	handler.MakeParticipantHandler(r, partSvc, eventSvc)
	handler.MakeUserHandler(r, userSvc)

	//_ = orgSvc.SaveOrg(&organization.Organization{
	//	Name:        "DSC VIT",
	//	Location:    "Vellore",
	//	Description: "The best",
	//	Tag:         "#dscvit",
	//	Website:     "dscvit.com",
	//	CreatedAt:   time.Now(),
	//})

	//_ = eventSvc.SaveEvent(&event.Event{
	//	ID:                    2,
	//	OrganizationID:        1,
	//	Name:                  "Leet haxx",
	//	Budget:                "690000",
	//	Description:           "Greatest hackathon ever",
	//	Category:              "Hackathon",
	//	Venue:                 "Some gallery",
	//	Attendance:            "full",
	//	ExpectedParticipants:  "5000",
	//	ToDate:                time.Now(),
	//	FromDate:              time.Now(),
	//	ToTime:                time.Now().Add(time.Hour),
	//	FromTime:              time.Now(),
	//})

	events, _ := orgSvc.GetOrgEvents(1)
	log.Println(events)

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

package main

import (
	"fmt"
	"github.com/ATechnoHazard/hades-2/api/handler"
	"github.com/ATechnoHazard/hades-2/pkg/coupon"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/hades-2/pkg/guest"
	"github.com/ATechnoHazard/hades-2/pkg/organization"
	"github.com/ATechnoHazard/hades-2/pkg/participant"
	"github.com/ATechnoHazard/hades-2/pkg/segment"
	"github.com/ATechnoHazard/hades-2/pkg/user"
	"github.com/ATechnoHazard/janus"
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
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	log.SetOutput(os.Stdout)
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func initNegroni() *negroni.Negroni {
	n := negroni.New()
	n.Use(negronilogrus.NewCustomMiddleware(log.DebugLevel, &log.JSONFormatter{PrettyPrint: true}, "API requests"))
	n.Use(negroni.NewRecovery())
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

	//db.AutoMigrate(&entities.Participant{}, &entities.Event{}, &entities.Organization{}, &entities.User{},
	//	entities.JoinRequest{}, &entities.Coupon{}, &entities.Guest{}, &entities.EventSegment{})
	return db
}

func main() {
	r := httprouter.New() // create a router
	n := initNegroni()    // init negroni middleware
	n.UseHandler(r)       // wrap router with negroni middleware
	db := connectDb()     // migrate and connect to db
	j, err := janus.NewJanusMiddleware(db)
	if err != nil {
		log.Panic(err)
	}

	// Create postgres repos for all entities
	partRepo := participant.NewPostgresRepo(db)
	eventRepo := event.NewPostgresRepo(db)
	orgRepo := organization.NewPostgresRepo(db)
	userRepo := user.NewPostgresRepo(db)
	guestRepo := guest.NewPostgresRepo(db)
	couponRepo := coupon.NewPostgresRepo(db)
	segmentRepo := segment.NewPostgresRepo(db)

	// Create services using previously generated repos
	partSvc := participant.NewParticipantService(partRepo)
	eventSvc := event.NewEventService(eventRepo)
	orgSvc := organization.NewOrganizationService(orgRepo)
	userSvc := user.NewUserService(userRepo)
	guestSvc := guest.NewGuestService(guestRepo)
	couponSvc := coupon.NewCouponService(couponRepo)
	segmentSvc := segment.NewEventSegmentService(segmentRepo)

	// Create and register handlers using generated services
	handler.MakeParticipantHandler(r, partSvc, eventSvc)
	handler.MakeUserHandler(r, userSvc)
	handler.MakeOrgHandler(r, orgSvc, j)
	handler.MakeGuestHandler(r, guestSvc, eventSvc)
	handler.MakeCouponHandler(r, couponSvc, eventSvc)
	handler.MakeEventSegmentHandler(r, segmentSvc, eventSvc)
	handler.MakeEventHandler(r, eventSvc, segmentSvc)
	handler.MakeAclHandler(r, j)

	// listen and serve on given port
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.WithField("event", "START").Info("Listening on port " + port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), n)
	if err != nil {
		log.Panic(err)
	}
}

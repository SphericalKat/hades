package event

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Service interface {
	SaveEvent(event *entities.Event) error
	ReadEvent(eventID uint) (*entities.Event, error)
	DeleteEvent(eventID uint) error
}

type eventSvc struct {
	repo Repository
}

func NewEventService(r Repository) Service {
	return &eventSvc{repo: r}
}

func (e *eventSvc) SaveEvent(event *entities.Event) error {
	return e.repo.Save(event)
}

func (e *eventSvc) ReadEvent(eventID uint) (*entities.Event, error) {
	return e.repo.Find(eventID)
}

func (e *eventSvc) DeleteEvent(eventID uint) error {
	return e.DeleteEvent(eventID)
}

package event

type Service interface {
	SaveEvent(event *Event) error
	ReadEvent(eventID uint) (*Event, error)
	DeleteEvent(eventID uint) error
}

type eventSvc struct {
	repo Repository
}

func NewEventService(r Repository) Service {
	return &eventSvc{repo: r}
}

func (e *eventSvc) SaveEvent(event *Event) error {
	return e.repo.Save(event)
}

func (e *eventSvc) ReadEvent(eventID uint) (*Event, error) {
	return e.repo.Find(eventID)
}

func (e *eventSvc) DeleteEvent(eventID uint) error {
	return e.DeleteEvent(eventID)
}

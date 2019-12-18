package event

type Repository interface {
	Find(eventID uint) (*Event, error)
	Save(event *Event) error
	Delete(eventID uint) error
}

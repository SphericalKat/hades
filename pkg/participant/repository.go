package participant

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	FindAll(eventId uint) ([]entities.Participant, error)
	FindByEmail(email string, eventID uint) (*entities.Participant, error)
	Save(participant *entities.Participant, eventID uint) error
	Delete(email string) error
	RemoveAttendeeEvent(email string, eventID uint) error
}

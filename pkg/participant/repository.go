package participant

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	FindAll() ([]entities.Participant, error)
	FindByRegNo(regNo string) (*entities.Participant, error)
	Save(participant *entities.Participant, eventID uint) error
	Delete(regNo string) error
	RemoveAttendeeEvent(regNo string, eventID uint) error
}

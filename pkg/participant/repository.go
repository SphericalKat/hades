package participant

type Repository interface {
	FindAll() ([]Participant, error)
	FindByRegNo(regNo string) (*Participant, error)
	Save(participant *Participant, eventID string) error
	Delete(regNo string) error
	RemoveAttendeeEvent(regNo string, eventID string) error
}

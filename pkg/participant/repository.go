package participant

type Repository interface {
	FindAll() ([]Participant, error)
	FindByRegNo(regNo string) (*Participant, error)
	Save(participant *Participant, eventID uint) error
	Delete(regNo string) error
	RemoveAttendeeEvent(regNo string, eventID uint) error
}

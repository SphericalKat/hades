package attendee

import "hades-2.0/pkg/entity"

type Reader interface {
	FindByRegNo(regNo string) (*entity.Attendee, error)
}

type Writer interface {
	Save(p *entity.Attendee) error
	Delete(regNo string) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Reader
	Writer
}

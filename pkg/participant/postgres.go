package participant

import (
	"github.com/jinzhu/gorm"
	"hades-2.0/pkg"
	"hades-2.0/pkg/event"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}

func (r *repo) FindAll() ([]Participant, error) {
	var participants []Participant
	err := r.DB.Model(Participant{}).Find(&participants).Error
	switch err {
	case nil:
		return participants, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) FindByRegNo(regNo string) (*Participant, error) {
	p := &Participant{}
	err := r.DB.Where("reg_no = ?", regNo).Find(p).Error
	switch err {
	case nil:
		return p, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) Save(participant *Participant, eventID string) error {
	tx := r.DB.Begin()
	e := &event.Event{EventID: eventID}

	if tx.Find(e).Error == gorm.ErrRecordNotFound {
		return pkg.ErrNotFound
	}

	err := tx.Save(participant).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Model(participant).Association("Events").Append(e).Error
	if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}
	tx.Commit()
	return nil
}

func (r *repo) Delete(regNo string) error {
	p := &Participant{}
	err := r.DB.Where("reg_no = ?", regNo).Delete(p).Error

	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) RemoveAttendeeEvent(regNo string, eventID string) error {
	tx := r.DB.Begin()
	e := &event.Event{EventID: eventID}
	p := &Participant{RegNo: regNo}

	if tx.Find(e).Error == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	}
	if tx.Find(p).Error == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	}

	err := tx.Model(p).Association("Events").Delete(e).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}
	tx.Commit()
	return nil
}

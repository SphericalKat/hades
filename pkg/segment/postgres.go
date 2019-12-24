package segment

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}

func (r *repo) GetSegments(eventId uint) ([]entities.EventSegment, error) {
	var segments []entities.EventSegment
	if err := r.DB.Where("event_id = ?", eventId).Find(&segments).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.ErrNotFound
		}
		return nil, pkg.ErrDatabase
	}

	return segments, nil
}

func (r *repo) DeleteSegment(segmentId uint) error {
	tx := r.DB.Begin()

	if err := tx.Delete(&entities.EventSegment{Day: segmentId}).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return pkg.ErrNotFound
		}
		return pkg.ErrDatabase
	}

	tx.Commit()
	return nil
}

func (r *repo) GetParticipantsInSegment(day uint) ([]entities.Participant, error) {
	eveSegment := &entities.EventSegment{Day: day}
	err := r.DB.Find(eveSegment).Association("PresentParticipants").Find(&eveSegment.PresentParticipants).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.ErrNotFound
		}
		return nil, pkg.ErrDatabase
	}

	return eveSegment.PresentParticipants, nil
}

func (r *repo) AddParticipantToSegment(regNo string, day uint, eventID uint) error {
	tx := r.DB.Begin()
	part := &entities.Participant{}
	eveSeg := &entities.EventSegment{}

	err := tx.Where("reg_no = ?", regNo).Find(part).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Where("day = ?", day).Where("event_id = ?", eventID).Find(eveSeg).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Find(eveSeg).Association("PresentParticipants").Find(&eveSeg.PresentParticipants).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	for _, p := range eveSeg.PresentParticipants {
		if p.RegNo == part.RegNo {
			tx.Rollback()
			return pkg.ErrAlreadyExists
		}
	}

	err = tx.Find(eveSeg).Association("PresentParticipants").Append(part).Error
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

func (r *repo) Find(day uint, eventID uint) (*entities.EventSegment, error) {
	seg := &entities.EventSegment{}
	tx := r.DB.Begin()
	err := tx.Where("day = ?", day).Where("event_id = ?", eventID).Find(seg).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, pkg.ErrNotFound
		default:
			return nil, pkg.ErrDatabase
		}
	}

	err = tx.Find(seg).Association("PresentParticipants").Find(&seg.PresentParticipants).Error
	switch err {
	case nil:
		tx.Commit()
		return seg, nil
	case gorm.ErrRecordNotFound:
		tx.Rollback()
		return nil, pkg.ErrNotFound
	default:
		tx.Rollback()
		return nil, pkg.ErrDatabase
	}
}

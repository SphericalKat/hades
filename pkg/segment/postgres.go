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

func (r *repo) GetParticipantsInSegment(segmentId uint) ([]entities.Participant, error) {
	eveSegment := &entities.EventSegment{Day: segmentId}
	err := r.DB.Find(eveSegment).Association("PresentParticipants").Find(&eveSegment.PresentParticipants).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.ErrNotFound
		}
		return nil, pkg.ErrDatabase
	}

	return eveSegment.PresentParticipants, nil
}

func (r *repo) AddPartipantToSegment(regNo string, day uint) error {
	tx := r.DB.Begin()
	part := &entities.Participant{}
	eveSeg := &entities.EventSegment{Day: day}

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

	err = tx.Where("day = ?", day).Find(eveSeg).Association("PresentParticipants").Append(part).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	return nil
}

func (r *repo) Find(day uint) (*entities.EventSegment, error) {
	seg := &entities.EventSegment{}
	if err := r.DB.Where("day = ?", day).Find(seg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.ErrNotFound
		}
		return nil, pkg.ErrDatabase
	}
	return seg, nil
}

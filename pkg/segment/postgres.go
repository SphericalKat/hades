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

func (r *repo) AddSegment(segment *entities.EventSegment) error {
	// Check for valid eventId.
	if err := r.DB.Find(&entities.Event{ID: segment.EventID}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return pkg.ErrNotFound
		}
		return pkg.ErrDatabase
	}

	if err := r.DB.Save(segment).Error; err != nil {
		return pkg.ErrDatabase
	}
	return nil
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

	if err := tx.Delete(&entities.EventSegment{SegmentID:segmentId}).Error; err != nil {
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
	var peeps []entities.Participant
	eveSegment := &entities.EventSegment{SegmentID:segmentId}
	err := r.DB.Find(eveSegment).Association("PresentParticipants").Find(&peeps).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.ErrNotFound
		}
		return nil, pkg.ErrDatabase
	}

	return peeps, nil
}

func (r *repo) AddPartipantToSegment(regNo string, segmentId uint) error {
	peep := &entities.Participant{RegNo:regNo}
	eveSeg := &entities.EventSegment{SegmentID:segmentId}

	if err := r.DB.Find(peep).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return pkg.ErrNotFound
		}
		return pkg.ErrDatabase
	}

	if err := r.DB.Find(eveSeg).Association("PresentParticipants").Append(peep).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return pkg.ErrNotFound
		}
		return pkg.ErrDatabase
	}

	return nil
}

func (r *repo) Find(segmentId uint) (*entities.EventSegment, error) {
	seg := &entities.EventSegment{SegmentID:segmentId}
	if err := r.DB.Find(seg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.ErrNotFound
		}
		return nil, pkg.ErrDatabase
	}
	return seg, nil
}

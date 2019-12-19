package guest

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB:db}
}

func (r *repo) Find(email string) (*Guest, error) {
	gst := &Guest{}
	err := r.DB.Where("email = ?", email).Find(gst).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	case nil:
		return gst, nil
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) FindAll() ([]Guest, error) {
	var guests []Guest
	err := r.DB.Model(Guest{}).Find(guests).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	return guests, nil
}

func (r *repo) Save(guest *Guest) error {
	tx := r.DB.Begin()
	err := tx.Save(guest).Error
	if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}
	tx.Commit()
	return nil
}

func (r *repo) Delete(email string) error {
	tx := r.DB.Begin()
	err := tx.Where("email = ?", email).Delete(&Guest{}).Error
	if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}
	tx.Commit()
	return nil
}



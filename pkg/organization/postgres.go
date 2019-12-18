package organization

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}

func (r *repo) Save(organization *Organization) error {
	err := r.DB.Save(organization).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) Find(orgID uint) (*Organization, error) {
	org := &Organization{ID: orgID}
	err := r.DB.Find(org).Association("Events").Find(&org.Events).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	case nil:
		return org, err
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) FindAll() ([]Organization, error) {
	var organizations []Organization
	err := r.DB.Model(Organization{}).Find(&organizations).Error
	return organizations, err
}

func (r *repo) Delete(orgID uint) error {
	err := r.DB.Delete(&Organization{ID: orgID}).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	case nil:
		return nil
	default:
		return pkg.ErrDatabase
	}
}

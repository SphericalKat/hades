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

func (rp *repo) Create(organization *Organization) error {
	_, err := rp.Find(organization.Name)
	switch err {
	case nil:
		return pkg.ErrAlreadyExists
	case pkg.ErrNotFound:
		err = rp.DB.Save(organization).Error
		if err != nil {
			return pkg.ErrDatabase
		}
		return err
	default:
		return pkg.ErrDatabase
	}
}

func (rp *repo) Find(name string) (*Organization, error) {
	org := &Organization{}
	err := rp.DB.Where("name = ?", name).Find(org).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	case nil:
		return org, err
	default:
		return nil, pkg.ErrDatabase
	}
}

func (rp *repo) FindAll() ([]Organization, error) {
	var organizations []Organization
	err := rp.DB.Model(Organization{}).Find(&organizations).Error
	return organizations, err
}

func (rp *repo) Delete(name string) error {
	err := rp.DB.Where("name = ?", name).Delete(&Organization{}).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	case nil:
		return nil
	default:
		return pkg.ErrDatabase
	}
}

func (rp *repo) Update(name string, organization *Organization) error {
	tx := rp.DB.Begin()

	err := tx.Save(organization).Error
	switch err {
	case gorm.ErrRecordNotFound:
		tx.Rollback()
		return pkg.ErrNotFound
	case nil:
		tx.Commit()
		return nil
	default:
		return pkg.ErrDatabase
	}
}
package organization

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func (r *repo) Save(organization *entities.Organization) (*entities.Organization, error) {
	err := r.DB.Save(organization).Error
	switch err {
	case nil:
		return organization, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) Find(orgID uint) (*entities.Organization, error) {
	org := &entities.Organization{ID: orgID}
	tx := r.DB.Begin()
	// find org
	err := tx.Where("id = ?", orgID).Find(org).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, pkg.ErrNotFound
		default:
			return nil, pkg.ErrDatabase
		}
	}

	// find org events
	err = tx.Find(org).Association("Events").Find(&org.Events).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, pkg.ErrNotFound
		default:
			return nil, pkg.ErrDatabase
		}
	}

	// find org users
	err = tx.Find(org).Association("Users").Find(&org.Users).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, pkg.ErrNotFound
		default:
			return nil, pkg.ErrDatabase
		}
	}

	// find org join requests
	err = tx.Find(org).Association("JoinRequests").Find(&org.JoinRequests).Error
	switch err {
	case nil:
		tx.Commit()
		return org, nil
	case gorm.ErrRecordNotFound:
		tx.Rollback()
		return nil, pkg.ErrNotFound
	default:
		tx.Rollback()
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) FindAll() ([]entities.Organization, error) {
	var organizations []entities.Organization
	err := r.DB.Model(entities.Organization{}).Find(&organizations).Error
	return organizations, err
}

func (r *repo) Delete(orgID uint) error {
	err := r.DB.Delete(&entities.Organization{ID: orgID}).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	case nil:
		return nil
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) SaveJoinReq(request *entities.JoinRequest) error {
	tx := r.DB.Begin()
	u := &entities.User{Email: request.Email}
	org := &entities.Organization{ID: request.OrganizationID}
	if err := tx.Where("email = ?", request.Email).Find(u).Error; err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	// find association
	if err := tx.Model(u).Association("Organizations").Find(&u.Organizations).Error; err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	// check if the join request already exists
	for _, org := range u.Organizations {
		if org.ID == request.OrganizationID {
			tx.Rollback()
			return pkg.ErrAlreadyExists
		}
	}

	// find the organization
	if tx.Where("id = ?", request.OrganizationID).Find(org).Error == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	}

	err := tx.Find(request).Error
	switch err {
	case gorm.ErrRecordNotFound:
		tx.Create(request)
		tx.Commit()
		return nil
	case nil:
		tx.Rollback()
		return pkg.ErrAlreadyExists
	default:
		tx.Rollback()
		return pkg.ErrDatabase
	}
}

func (r *repo) DelJoinReq(orgID uint, email string) error {
	org := &entities.Organization{}
	tx := r.DB.Begin()
	if email == "" || orgID == 0 {
		return pkg.ErrInvalidSlug
	}

	// find organization to delete the join req from
	if err := tx.Where("id = ?", orgID).Find(org).Error; err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	// find join requests belonging to the org
	if err := tx.Model(org).Association("JoinRequests").Delete(&entities.JoinRequest{
		Email: email,
	}).Error; err != nil {
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

func (r *repo) FindAllJoinReq(orgID uint) ([]entities.JoinRequest, error) {
	org := &entities.Organization{}

	tx := r.DB.Begin()
	err := tx.Where("id = ?", orgID).Find(org).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, pkg.ErrNotFound
		default:
			return nil, pkg.ErrDatabase
		}
	}
	err = tx.Find(org).Association("JoinRequests").Find(&org.JoinRequests).Error
	switch err {
	case nil:
		tx.Commit()
		return org.JoinRequests, nil
	case gorm.ErrRecordNotFound:
		tx.Rollback()
		return nil, pkg.ErrNotFound
	default:
		tx.Rollback()
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) AcceptJoinReq(request *entities.JoinRequest) error {
	tx := r.DB.Begin()
	u := &entities.User{Email: request.Email}
	org := &entities.Organization{ID: request.OrganizationID}
	err := tx.Find(request).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Find(org).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Find(u).Association("Organizations").Append(org).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrAlreadyExists
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Delete(request).Error
	switch err {
	case nil:
		tx.Commit()
		return nil
	case gorm.ErrRecordNotFound:
		tx.Rollback()
		return pkg.ErrAlreadyExists
	default:
		tx.Rollback()
		return pkg.ErrDatabase
	}
}

func (r *repo) AddUserToOrg(orgID uint, email string) error {
	tx := r.DB.Begin()
	org := &entities.Organization{}
	u := &entities.User{}

	// find if org exists
	err := tx.Where("id = ?", orgID).Find(org).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrAlreadyExists
		default:
			return pkg.ErrDatabase
		}
	}

	// find user
	err = tx.Where("email = ?", email).Find(u).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrAlreadyExists
		default:
			return pkg.ErrDatabase
		}
	}

	// add user to org
	err = tx.Find(org).Association("Users").Append(u).Error
	if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}
	tx.Commit()
	return nil
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}

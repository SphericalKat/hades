package attendee

import (
	"github.com/jinzhu/gorm"
	"hades-2.0/pkg/entity"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (pr *PostgresRepository) FindByRegNo(regNo string) (*entity.Attendee, error) {
	p := &entity.Attendee{}
	err := pr.db.Where("registration_number = ?", regNo).Find(p).Error
	return p, err
}

func (pr *PostgresRepository) Save(p *entity.Attendee) error {
	return pr.db.Save(p).Error
}

func (pr *PostgresRepository) Delete(regNo string) error {
	p := &entity.Attendee{RegistrationNumber: regNo}
	return pr.db.Delete(p).Error
}

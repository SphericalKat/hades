package organization

import (
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"time"
)

type Organization struct {
	ID          uint          `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name        string        `json:"name"`
	Location    string        `json:"location"`
	Description string        `json:"description"`
	Tag         string        `json:"tag"`
	Website     string        `json:"website"`
	CreatedAt   time.Time     `json:"created_at"`
	DeletedAt   *time.Time    `sql:"index"`
	Events      []event.Event `json:"-"`
}

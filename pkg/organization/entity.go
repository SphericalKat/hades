package organization

import "time"

type Organization struct {
	Name        string    `json:"name" gorm:"primary_key;"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	Website     string    `json:"website"`
	CreatedAt   time.Time `json:"createdAt"`
}

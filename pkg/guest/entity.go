package guest

type Guest struct {
	Name           string `json:"name"`
	Email          string `json:"email" gorm:"primary_key;"`
	PhoneNumber    string `json:"phone_number"`
	Gender         string `json:"gender"`
	Stake          string `json:"stake"`
	LocationOfStay string `json:"location_of_stay"`
}

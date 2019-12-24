package entities

type CSVParticipant struct {
	Name        string `csv:"name"`
	RegNo       string `csv:"reg_no" gorm:"primary_key"`
	Email       string `csv:"email"`
	PhoneNumber string `csv:"phone_number"`
	Gender      string `csv:"gender"`
	IsPresent   bool   `csv:"is_present"`
}

func P2CSVPTransform(part *Participant) *CSVParticipant {
	return &CSVParticipant{
		Name:        part.Name,
		RegNo:       part.RegNo,
		Email:       part.Email,
		PhoneNumber: part.PhoneNumber,
		Gender:      part.Gender,
		IsPresent:   false,
	}
}

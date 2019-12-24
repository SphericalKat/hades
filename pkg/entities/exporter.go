package entities

type CSVParticipant struct {
	Name        string `csv:"name" json:"name"`
	RegNo       string `csv:"reg_no" json:"reg_no"`
	Email       string `csv:"email" json:"email"`
	PhoneNumber string `csv:"phone_number" json:"phone_number"`
	Gender      string `csv:"gender" json:"gender"`
	IsPresent   bool   `csv:"is_present" json:"is_present"`
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

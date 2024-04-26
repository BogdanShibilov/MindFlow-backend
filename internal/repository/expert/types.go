package expertrepo

type ProffFieldListAndCount struct {
	ProfessionalField string `db:"professional_field" json:"professionalField"`
	Count             int    `db:"count" json:"count"`
}

type MinMaxPrice struct {
	MinPrice int `json:"minPrice" db:"min_price"`
	MaxPrice int `json:"maxPrice" db:"max_price"`
}

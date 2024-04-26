package expertservice

import expertrepo "github.com/bogdanshibilov/mindflowbackend/internal/repository/expert"

type FilterData struct {
	ProffFieldData []expertrepo.ProffFieldListAndCount `json:"proffFieldData"`
	MinMaxPrice    expertrepo.MinMaxPrice              `json:"minMaxPrice"`
}

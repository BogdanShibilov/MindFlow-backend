package entity

type User struct {
	Uuid     string
	Email    string
	PassHash []byte
}

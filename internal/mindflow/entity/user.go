package entity

type User struct {
	Guid     string
	Email    string
	PassHash []byte
}

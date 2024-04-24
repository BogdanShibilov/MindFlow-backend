package entity

type Status int

const (
	Pending Status = iota
	Approved
	Rejected
)

type Permission int

const (
	Admin Permission = iota
)

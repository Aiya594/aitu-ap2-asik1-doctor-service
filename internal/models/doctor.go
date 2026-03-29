package models

type Doctor struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization"`
}

package model

import "time"

type DoctorCreated struct {
	EventType      string    `json:"event_type"`
	OccurredAt     time.Time `json:"occurred_at"`
	ID             string    `json:"id"`
	Full_name      string    `json:"full_name"`
	Specialization string    `json:"specialization"`
	Email          string    `json:"email"`
}

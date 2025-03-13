package models

import "time"

// fault notif struct to represent what we get from validator and what we will input into firestore notifications doc
type FaultNotification struct {
	ProjectID string    `firestore:"project_id" json:"project_id"`
	Message   string    `firestore:"message" json:"message"`
	StartTime time.Time `firestore:"start_time" json:"start_time"`
	EndTime   time.Time `firestore:"end_time" json:"end_time"`
	Average   float64   `firestore:"average" json:"average"`
}

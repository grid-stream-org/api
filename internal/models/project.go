package models

// defines project model

import "time"

// Represents an offloading instance from a DER(s)
type Project struct {
	ProjectID         string    `json:"projectId"`         // Unique identifier for the project
	UtilityID         string    `json:"utilityId"`         // Unique identifier for the utility
	ConnectionStartAt time.Time `json:"connectionStartAt"` // Timestamp when the connection started
}

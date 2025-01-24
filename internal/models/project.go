package models

// Represents an offloading instance from a DER(s)
type Project struct {
	ID        string `json:"id" bigquery:"id"`                 // Unique identifier for the project
	UtilityID string `json:"utility_id" bigquery:"utility_id"` // Unique identifier for the utility
	UserID    string `json:"user_id" bigquery:"user_id"`
	Location  string `json:"location" bigquery:"location"`
}

package models

// defines project model

// Represents an offloading instance from a DER(s)
type Project struct {
	Id        string `json:"id"`         // Unique identifier for the project
	UtilityId string `json:"utility_id"` // Unique identifier for the utility
	UserId    string `json:"user_id"`
	Location  string `json:"location"`
}

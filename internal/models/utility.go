package models

type Utility struct {
	ID          string `json:"id" bigquery:"id"`
	DisplayName string `json:"display_name" bigquery:"display_name"`
}

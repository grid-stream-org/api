package models

type DERMetadata struct {
	ID                string  `json:"id" bigquery:"id"`
	ProjectID         string  `json:"project_id" bigquery:"project_id"`
	Type              string  `json:"type" bigquery:"type"`
	NameplateCapacity float64 `json:"nameplate_capacity" bigquery:"nameplate_capacity"`
	PowerCapacity     float64 `json:"power_capacity" bigquery:"power_capacity"`
}

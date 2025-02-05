package models

type DERType string
type DERMetadata struct {
	ID                string  `json:"id" bigquery:"id"`
	ProjectID         string  `json:"project_id" bigquery:"project_id"`
	Type              DERType `json:"type" bigquery:"type"`
	NameplateCapacity float64 `json:"nameplate_capacity" bigquery:"nameplate_capacity"`
	PowerCapacity     float64 `json:"power_capacity" bigquery:"power_capacity"`
}

const (
	Solar   DERType = "solar"
	Battery DERType = "battery"
	EV      DERType = "ev"
)


func (t DERType) IsValid() bool {
	switch t {
	case Solar, Battery, EV:
		return true
	default:
		return false
	}
}
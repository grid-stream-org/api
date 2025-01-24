package models

type DERMetadata struct {
	ID                string  `json:"id"`                
	Type              string  `json:"type"`               
	NameplateCapacity float64 `json:"nameplate_capacity"` 
	PowerCapacity     float64 `json:"power_capacity"`     
}

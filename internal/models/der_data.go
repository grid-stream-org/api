package models

import "time"

type DERData struct {
	ID                    string    `json:"id"`
	DERID                 string    `json:"der_id"`
	Timestamp             time.Time `json:"timestamp"`
	CurrentOutput         float64   `json:"current_output"`
	Units                 string    `json:"units"`
	ProjectID             string    `json:"project_id"`
	Baseline              float64    `json:"baseline"`
	IsOnline              bool      `json:"is_online"`
	IsStandalone          bool      `json:"is_standalone"`
	ConnectionStartAt     string    `json:"connection_start_at"`
	CurrentSOC            float64   `json:"current_soc"`
	PowerMeterMeasurement float64   `json:"power_meter_measurement"`
	ContractThreshold     float64   `json:"contract_threshold"`
}

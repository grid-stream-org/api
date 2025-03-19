package models

import "time"

// ProjectAverage represents a project's average output data
type ProjectAverage struct {
	ProjectID         string    `json:"project_id" bigquery:"project_id"`
	StartTime         time.Time `json:"start_time" bigquery:"start_time"`
	EndTime           time.Time `json:"end_time" bigquery:"end_time"`
	Baseline          float64   `json:"baseline" bigquery:"baseline"`
	ContractThreshold float64   `json:"contract_threshold" bigquery:"contract_threshold"`
	AverageOutput     float64   `json:"average_output" bigquery:"average_output"`
}
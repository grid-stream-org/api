package models

import "time"

type ProjectAverage struct {
	Timestamp     time.Time `json:"timestamp"`
	ProjectID     string    `json:"project_id"`
	AverageOutput float64   `json:"average_output"`
}

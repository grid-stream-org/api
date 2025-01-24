package models

type ProjectAverage struct {
    Timestamp     string  `json:"timestamp"`      
    ProjectID     string  `json:"project_id"`     
    AverageOutput float64 `json:"average_output"` 
}

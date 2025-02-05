package models

import "cloud.google.com/go/bigquery"

type ContractStatus string
type Contract struct {
	ID                string            `json:"id" bigquery:"id"`
	ContractThreshold float64           `json:"contract_threshold" bigquery:"contract_threshold"`
	StartDate         bigquery.NullDate `json:"start_date" bigquery:"start_date"` // 2006-01-02T15:04:05Z RFC-3339 format
	EndDate           bigquery.NullDate `json:"end_date" bigquery:"end_date"`
	Status            ContractStatus    `json:"status" bigquery:"status"`
	ProjectID         string            `json:"project_id" bigquery:"project_id"`
}

const (
	Active   ContractStatus = "active"
	Inactive ContractStatus = "inactive"
	Pending  ContractStatus = "pending"
)

func (s ContractStatus) IsValid() bool {
	switch s {
	case Active, Pending, Inactive:
		return true
	default:
		return false
	}
}

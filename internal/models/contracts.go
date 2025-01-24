package models

type Contract struct {
	ID                string  `json:"id"`
	ContractThreshold float64 `json:"contract_threshold"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date"`
	Status            string  `json:"status"`
	ProjectID         string  `json:"project_id"`
}

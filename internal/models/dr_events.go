package models

import (
	"time"
)

type DREvents struct {
	ID        string    `json:"id" bigquery:"id"`
	UtilityID string    `json:"utility_id" bigquery:"utility_id"`
	StartTime time.Time `json:"start_time" bigquery:"start_time"`
	EndTime   time.Time `json:"end_time" bigquery:"end_time"`
}

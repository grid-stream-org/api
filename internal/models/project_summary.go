package models

import "time"

type ProjectSummary struct {
	TotalActive       int       `json:"total_active" bigquery:"total_active"`
	TotalPending      int       `json:"total_pending" bigquery:"total_pending"`
	TotalThreshold    float64       `json:"total_threshold" bigquery:"total_threshold"`
	NextEventID       string    `json:"next_event_id" bigquery:"next_event_id"`
	NextEventStart    time.Time `json:"next_event_start" bigquery:"next_event_start"`
	NextEventEnd      time.Time `json:"next_event_end" bigquery:"next_event_end"`
	RecentEventID     string    `json:"recent_event_id" bigquery:"recent_event_id"`
	RecentEventStart  time.Time `json:"recent_event_start" bigquery:"recent_event_start"`
	RecentEventEnd    time.Time `json:"recent_event_end" bigquery:"recent_event_end"`
}
package model

import "time"

type GenderizeResponse struct {
	Count       int      `json:"count"`
	Name        string   `json:"name"`
	Gender      *string  `json:"gender"`
	Probability float64  `json:"probability"`
}

type ClassifyResponseData struct {
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	Probability float64   `json:"probability"`
	SampleSize  int       `json:"sample_size"`
	IsConfident bool      `json:"is_confident"`
	ProcessedAt time.Time `json:"processed_at"`
}

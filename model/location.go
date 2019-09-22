package model

import "time"

// Location represents a location sent by the phone
type Location struct {
	ID             uint       `gorm:"primary_key" json:"-"`
	CreatedAt      time.Time  `json:"at"`
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `sql:"index" json:"-"`
	CoarseLocation string     `json:"coarse_location"`
	SearchQuery    string     `json:"search_query"`
	Latitude       float64    `json:"latitude,omitempty"`
	Longitude      float64    `json:"longitude,omitempty"`
	Accuracy       float32    `json:"accuracy,omitempty"`
	Battery        float32    `json:"battery"`
}

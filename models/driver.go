package models

import "time"

type Driver struct {
	ID              string     `db:"id" json:"id"`
	UserID          string     `db:"user_id" json:"userId"`
	IsOnline        bool       `db:"is_online" json:"isOnline"`
	IsAvailable     bool       `db:"is_available" json:"isAvailable"`
	CurrentLatitude *float64   `db:"current_latitude" json:"currentLatitude"`
	CurrentLongitude *float64  `db:"current_longitude" json:"currentLongitude"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updatedAt"`
	DLnce float64 `db:"distance" json:"distance,omitempty"`
}


package models

import "time"

type Driver struct {
	ID              string     `db:"id" json:"id"`
	UserID          string     `db:"user_id" json:"user_id"`
	IsOnline        bool       `db:"is_online" json:"is_online"`
	IsAvailable     bool       `db:"is_available" json:"is_available"`
	CurrentLatitude *float64   `db:"current_latitude" json:"current_latitude"`
	CurrentLongitude *float64  `db:"current_longitude" json:"current_longitude"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}


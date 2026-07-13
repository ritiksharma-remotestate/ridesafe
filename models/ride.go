package models

import "time"

type Ride struct {
	ID string `db:"id" json:"id"`
	PassengerID string `db:"passenger_id" json:"passenger_id"`
	DriverID *string `db:"driver_id" json:"driver_id"`
	PickupLatitude  float64 `db:"pickup_latitude" json:"pickup_latitude"`
	PickupLongitude float64 `db:"pickup_longitude" json:"pickup_longitude"`
	DestinationLatitude  float64 `db:"destination_latitude" json:"destination_latitude"`
	DestinationLongitude float64 `db:"destination_longitude" json:"destination_longitude"`
	Fare *float64 `db:"fare" json:"fare"`
	Status string `db:"status" json:"status"`
	RequestedAt time.Time `db:"requested_at" json:"requested_at"`
	AcceptedAt  *time.Time `db:"accepted_at" json:"accepted_at"`
	ArrivedAt   *time.Time `db:"arrived_at" json:"arrived_at"`
	StartedAt   *time.Time `db:"started_at" json:"started_at"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at"`
	CancelledAt *time.Time `db:"cancelled_at" json:"cancelled_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
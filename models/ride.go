package models

import "time"

type RideStatus string

const (
	RideRequested RideStatus = "REQUESTED"
	RideAccepted  RideStatus = "ACCEPTED"
	RideArrived   RideStatus = "DRIVER_ARRIVED"
	RideStarted   RideStatus = "ONGOING"
	RideCompleted RideStatus = "COMPLETED"
	RideCancelled RideStatus = "CANCELLED"
)

type Ride struct {
	ID                   string     `db:"id" json:"id"`
	PassengerID          string     `db:"passenger_id" json:"passengerId"`
	DriverID             *string    `db:"driver_id" json:"driverId"`
	PickupLatitude       float64    `db:"pickup_latitude" json:"pickupLatitude"`
	PickupLongitude      float64    `db:"pickup_longitude" json:"pickupLongitude"`
	DestinationLatitude  float64    `db:"destination_latitude" json:"destinationLatitude"`
	DestinationLongitude float64    `db:"destination_longitude" json:"destinationLongitude"`
	Fare                 *float64   `db:"fare" json:"fare"`
	Status               string     `db:"status" json:"status"`
	RequestedAt          time.Time  `db:"requested_at" json:"requestedAt"`
	AcceptedAt           *time.Time `db:"accepted_at" json:"acceptedAt"`
	ArrivedAt            *time.Time `db:"arrived_at" json:"arrivedAt"`
	StartedAt            *time.Time `db:"started_at" json:"startedAt"`
	CompletedAt          *time.Time `db:"completed_at" json:"completedAt"`
	CancelledAt          *time.Time `db:"cancelled_at" json:"cancelledAt"`
	CreatedAt            time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time  `db:"updated_at" json:"updatedAt"`
	RideOTPHash          *string    `db:"ride_otp_hash"`
	RideOTPGeneratedAt   *time.Time `db:"ride_otp_generated_at"`
	OTPVerifiedAt        *time.Time `db:"otp_verified_at"`
	RideOTPVerified      bool       `db:"ride_otp_verified" json:"RideOTPVerified"`
	StartOTP             *string    `db:"started_otp" json:"-"`
}

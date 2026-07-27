package models

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     Role   `json:"role" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type UpdateDriverLocationRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,gte=-90,lte=90"`
	Longitude float64 `json:"longitude" validate:"required,gte=-180,lte=180"`
}
type DriverOnlineRequest struct {
	Online bool `json:"online"`
}

type DriverAvailableRequest struct {
	Available bool `json:"available"`
}
type CreateRideRequest struct {
	PickupLatitude       float64 `json:"pickupLatitude" validate:"required"`
	PickupLongitude      float64 `json:"pickupLongitude" validate:"required"`
	DestinationLatitude  float64 `json:"destinationLatitude" validate:"required"`
	DestinationLongitude float64 `json:"destinationLongitude" validate:"required"`
}
type AvailableDrivers struct {
	PickupLatitude  float64 `json:"pickupLatitude" validate:"required"`
	PickupLongitude float64 `json:"pickupLongitude" validate:"required"`
}
type VerifyOTPRequest struct {
	OTP string `json:"otp" validate:"required,len=4,numeric"`
}

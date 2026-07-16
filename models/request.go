package models




type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role 	Role 	`json:"role"`
	
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	PickupLatitude       float64 `json:"pickup_latitude" validate:"required"`
	PickupLongitude      float64 `json:"pickup_longitude" validate:"required"`
	DestinationLatitude  float64 `json:"destination_latitude" validate:"required"`
	DestinationLongitude float64 `json:"destination_longitude" validate:"required"`
}
type AvailableDrivers struct{
	PickupLatitude       float64 `json:"pickup_latitude" validate:"required"`
	PickupLongitude      float64 `json:"pickup_longitude" validate:"required"`
}
type VerifyOTPRequest struct {
	OTP string `json:"otp" validate:"required,len=4,numeric"`
}
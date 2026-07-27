package error_custom

import "errors"

// User Errors
var (
	ErrUserNotFound        = errors.New("user not found")
	ErrDriverNotFound      = errors.New("driver not found")
	ErrPassengerNotFound   = errors.New("passenger not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidRole         = errors.New("invalid role")
	ErrNotADriver          = errors.New("user is not a driver")
	ErrNotAPassenger       = errors.New("user is not a passenger")
	ErrDriverAlreadyExists = errors.New("driver already exists")
)

// Ride Error
var (
	ErrOTPNotVerified       = errors.New("otp not verified")
	ErrInvalidOTP           = errors.New("invalid otp")
	ErrRideNotFound         = errors.New("ride not found")
	ErrRideAlreadyAccepted  = errors.New("ride has already been accepted")
	ErrRideAlreadyStarted   = errors.New("ride has already started")
	ErrRideAlreadyCompleted = errors.New("ride has already been completed")
	ErrRideAlreadyCancelled = errors.New("ride has already been cancelled")

	ErrRideNotAccepted  = errors.New("ride has not been accepted")
	ErrRideNotStarted   = errors.New("ride has not been started")
	ErrRideNotCompleted = errors.New("ride has not been completed")

	ErrTooManyOTPAttempts    = errors.New("too many otp attempts")
	ErrInvalidRideStatus     = errors.New("invalid ride status")
	ErrInvalidRideTransition = errors.New("invalid ride status transition")
)

// Driver Errors
var (
	ErrDriverUnavailable     = errors.New("driver is unavailable")
	ErrDriverOffline         = errors.New("driver is offline")
	ErrDriverBusy            = errors.New("driver is already on another ride")
	ErrUnauthorizedDriver    = errors.New("unauthorized driver")
	ErrDriverLocationMissing = errors.New("driver location not available")
)

// Passenger Errors
var (
	ErrPassengerAlreadyOnRide = errors.New("passenger already has an active ride")
)

// Location Error
var (
	ErrLocationRequired    = errors.New("location is required")
	ErrInvalidLocation     = errors.New("invalid location")
	ErrDestinationRequired = errors.New("destination is required")
)

// Authentication Errors
var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionNotFound = errors.New("session not found")
)

// Validation Errors
var (
	ErrInvalidRequest        = errors.New("invalid request")
	ErrMissingRequiredFields = errors.New("missing required fields")
)

// Database Errors
var (
	ErrDatabase = errors.New("database error")
)

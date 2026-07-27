package services

import (
	"database/sql"
	"errors"
	"fmt"

	"ridesafe/database"
	"ridesafe/error_custom"
	"ridesafe/models"
	"ridesafe/repository"
	"ridesafe/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	tx *sqlx.DB
)

func CreateRide(passengerID string, req models.CreateRideRequest) (*models.Ride, error) {

	user, err := repository.GetUserByID(passengerID)
	if err != nil {
		return nil, error_custom.ErrUserNotFound
	}

	if user.Role != models.RolePassenger {
		return nil, error_custom.ErrNotAPassenger
	}

	distance := utils.CalculateDistance(req.PickupLatitude,req.PickupLongitude,req.DestinationLatitude,req.DestinationLongitude,)

	fare := utils.CalculateFare(distance)

	ride := models.Ride{
		PassengerID:          passengerID,
		PickupLatitude:       req.PickupLatitude,
		PickupLongitude:      req.PickupLongitude,
		DestinationLatitude:  req.DestinationLatitude,
		DestinationLongitude: req.DestinationLongitude,
		Fare:        &fare,
		Status:      string(models.RideRequested),
		RequestedAt: time.Now(),
	}

	newRide, err := repository.CreateRide(ride)
	if err != nil {
		return nil, err
	}

	return newRide, nil
}

func AcceptRide(rideID string, driverID string) error {

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		fmt.Print("60")
		return error_custom.ErrRideNotFound
	}
	if models.RideStatus(ride.Status) != models.RideRequested {
		return error_custom.ErrDriverAlreadyExists
	}

	driver, err := repository.GetDriverByUserID(driverID)
	if err != nil {
		return error_custom.ErrNotADriver
	}

	if !driver.IsOnline {
		return error_custom.ErrDriverOffline
	}
	if !driver.IsAvailable {
		return error_custom.ErrDriverUnavailable
	}

	return database.Tx(func(tx *sqlx.Tx) error {

		err := repository.AcceptRide(tx,rideID,driverID,time.Now(),models.RideAccepted,)
		if err != nil {
			return err
		}

		err = repository.SetDriverAvailable(tx,driverID,false,)
		if err != nil {
			return err
		}

		return nil
	})
}

func StartRide(rideID, driverID string) error {

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_custom.ErrRideNotFound
		}
		return err
	}

	if models.RideStatus(ride.Status) != models.RideArrived {
		return error_custom.ErrRideNotAccepted
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return error_custom.ErrUnauthorizedDriver
	}
	if ride.OTPVerifiedAt == nil {
		return error_custom.ErrOTPNotVerified
	}

	err = repository.StartRide(rideID,time.Now(),models.RideStarted,)
	
	if err != nil {
		return err
	}

	return nil

}

func CompleteRide(rideID, driverID string) error {

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_custom.ErrRideNotFound
		}
		return err
	}
	if models.RideStatus(ride.Status) == models.RideCompleted {
		return error_custom.ErrRideAlreadyCompleted
	}

	if models.RideStatus(ride.Status) != models.RideStarted {
		return error_custom.ErrRideNotStarted
	}


	if ride.DriverID == nil || *ride.DriverID != driverID {
		return error_custom.ErrUnauthorizedDriver
	}
	return database.Tx(func(tx *sqlx.Tx) error {

		err = repository.CompleteRide(tx,rideID,time.Now(),models.RideCompleted,)
		if err != nil {
			return err
		}

		err = repository.SetDriverAvailable(tx, driverID, true)
		if err != nil {
			return err
		}

		return nil
	})
}
func CancelRide(rideID string, userID string, role models.Role) error {

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_custom.ErrRideNotFound
		}
		return err
	}
	if role == models.RolePassenger {
		if ride.PassengerID != userID {
			return error_custom.ErrUnauthorized
		}
	}

	if role == models.RoleDriver {
		if ride.DriverID == nil || *ride.DriverID != userID {
			return error_custom.ErrUnauthorized
		}
	}

	if models.RideStatus(ride.Status) == models.RideCompleted {
		return error_custom.ErrRideAlreadyCompleted
	}

	if models.RideStatus(ride.Status) == models.RideCancelled {
		return error_custom.ErrRideAlreadyCancelled
	}
	return database.Tx(func(tx *sqlx.Tx) error {
		err = repository.CancelRide(tx,rideID,time.Now(),models.RideCancelled,)
		if err != nil {
			return err
		}

		if ride.DriverID != nil {
			err = repository.SetDriverAvailable(tx, *ride.DriverID, true)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func MarkRideArrived(rideID, driverID string) (string,error) {

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "",error_custom.ErrRideNotFound
		}
		return "",err
	}

	if models.RideStatus(ride.Status) != models.RideAccepted {
		return "",error_custom.ErrRideNotAccepted
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return "",error_custom.ErrUnauthorizedDriver
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return "",err
	}

	err = repository.SaveRideOTP(rideID,otp,time.Now(),models.RideArrived,)
	if err != nil {
		return "",err
	}

	return otp,nil
}
func VerifyRideOTP(rideID string,driverID string,otp string) error {

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_custom.ErrRideNotFound
		}
		return err
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return error_custom.ErrUnauthorizedDriver
	}

	if models.RideStatus(ride.Status) != models.RideArrived {
		return error_custom.ErrRideNotAccepted
	}

	err = repository.VerifyRideOTP(rideID,otp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_custom.ErrInvalidOTP
		}
		return err
	}

	return nil
}

func GetRideByID(rideID string, userID string, role models.Role) (*models.Ride,error){

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil,error_custom.ErrRideNotFound
		}
		return nil,err
	}
	if role == models.RolePassenger {
		if ride.PassengerID != userID {
			return nil,error_custom.ErrUnauthorized
		}
	}
	
	if role == models.RoleDriver {
		if ride.DriverID == nil || *ride.DriverID != userID {
			return nil,error_custom.ErrUnauthorized
		}
	}

	return ride,nil
	
}

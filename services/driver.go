package services

import (
	"database/sql"
	"errors"
	"ridesafe/database"
	"ridesafe/error_custom"
	"ridesafe/models"
	"ridesafe/repository"
	"ridesafe/utils"
)

func CreateDriver(user *utils.Claims) error {
	if models.Role(user.Role) != models.RoleDriver {
		return error_custom.ErrNotADriver
	}

	driver, err := repository.GetDriverByUserID(user.UserID)
	if err == nil && driver != nil {
		return error_custom.ErrDriverAlreadyExists
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = repository.CreateDriver(user.UserID)
	if err != nil {
		return err
	}

	return nil
}

func GetDriverByUserID(userID string) (*models.Driver, error) {
	driver, err := repository.GetDriverByUserID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error_custom.ErrUserNotFound
		}
		return nil, err

	}
	return driver, nil

}

func UpdateDriverLocation(userID string, latitude, longitude float64) error {
	_, err := repository.GetDriverByUserID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_custom.ErrUserNotFound
		}
		return err

	}

	err = repository.UpdateDriverLocation(userID, latitude, longitude)
	if err != nil {
		return err
	}
	return nil

}
func DriverOnline(userID string, online bool) error {
	_, err := repository.GetDriverByUserID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_custom.ErrUserNotFound
		}
		return err
	}
	err = repository.SetDriverOnline(userID, online)
	if err != nil {
		return err
	}
	return nil

}
func DriverAvailable(userID string, available bool) error {
	_, err := repository.GetDriverByUserID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_custom.ErrUserNotFound
		}
		return err
	}

	err = repository.SetDriverAvailable(database.Ridesafe, userID, available)
	if err != nil {
		return err
	}
	return nil

}

func GetAvailableDrivers(rideID, userID string) ([]models.Driver, error) {

	rideInfo, err := repository.GetRideByID(rideID)
	if err != nil {
		return nil, err
	}
	if rideInfo.PassengerID != userID {
		return nil, error_custom.ErrUnauthorized
	}
	latitude := rideInfo.PickupLatitude
	longitude := rideInfo.PickupLongitude

	drivers, err := repository.GetAvailableDrivers(latitude, longitude)
	if err != nil {
		return nil, err
	}
	return drivers, nil
}
func GetDriverLocation(userID string) (float64, float64, error) {
	lat, lon, err := repository.GetDriverLocation(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, error_custom.ErrDriverLocationMissing
		}
		return 0, 0, err
	}

	return lat, lon, nil
}

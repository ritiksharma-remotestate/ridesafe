package repository

import (
	"database/sql"
	"ridesafe/database"
	"ridesafe/models"
	"time"
	"github.com/jmoiron/sqlx"
)

func CreateRide(ride models.Ride) (*models.Ride, error) {
	SQL := `INSERT INTO rides ( passenger_id, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude, fare, status
)
VALUES ( $1, $2, $3, $4, $5, $6, $7)
RETURNING  id, passenger_id, driver_id, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude, fare, status, requested_at, accepted_at, arrived_at, started_at, completed_at, cancelled_at, ride_otp_hash, ride_otp_generated_at, created_at, updated_at;`

	var newRide models.Ride

	err := database.Ridesafe.Get(&newRide, SQL, ride.PassengerID, ride.PickupLatitude, ride.PickupLongitude, ride.DestinationLatitude, ride.DestinationLongitude, ride.Fare, ride.Status)
	if err != nil {
		return nil, err
	}
	return &newRide, nil
}

func AcceptRide(db sqlx.Ext, rideID string, driverID string, acceptedAt time.Time, status models.RideStatus) error {
	SQL := `UPDATE rides SET driver_id = $1, accepted_at = $2, status = $3 WHERE id = $4;`

	_, err := db.Exec(SQL, driverID, acceptedAt, status, rideID)
	if err != nil {
		return err
	}
	return nil
}

func StartRide(rideID string, started_at time.Time, status models.RideStatus) error {
	SQL := `UPDATE rides SET started_at= $1, status = $2 WHERE id = $3;`
	_, err := database.Ridesafe.Exec(SQL, started_at,status,rideID)
	if err != nil {
		return err
	}
	return nil
}

func CompleteRide( db sqlx.Ext, rideID string, completed_at time.Time, status models.RideStatus,) error {
	SQL := `UPDATE rides SET completed_at= $1, status = $2 WHERE id = $3;`

	_, err := db.Exec( SQL, completed_at, status, rideID,)
	if err != nil {
		return err
	}
	return nil
}
func CancelRide( db sqlx.Ext, rideID string, cancelled_at time.Time, status models.RideStatus,) error {
	SQL := `UPDATE rides SET cancelled_at= $1, status = $2 WHERE id = $3;`
	_, err := db.Exec( SQL, cancelled_at, status, rideID,)
	if err != nil {
		return err
	}
	return nil
}

func MarkRideArrived(rideID string, arrived_at time.Time, status models.RideStatus) error {
	SQL := `UPDATE rides SET arrived_at= $1,status = $2WHERE id = $3;`

	_, err := database.Ridesafe.Exec(SQL, arrived_at, status, rideID)
	if err != nil {
		return err
	}
	return nil
}

func GetRideByID(rideID string) (*models.Ride, error) {
	query := `
		SELECT  id, passenger_id, driver_id, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude, fare, status, requested_at, otp_verified_at, accepted_at, arrived_at, started_at, completed_at, cancelled_at, created_at, ride_otp_hash, ride_otp_generated_at, start_otp, updated_at FROM rides WHERE id = $1`

	var ride models.Ride

	err := database.Ridesafe.Get(&ride, query, rideID)

	if err != nil {
		return nil, err
	}

	return &ride, nil
}

func SaveRideOTP(rideID string, otp string, arrivedAt time.Time, status models.RideStatus) error {

	SQL := `
	UPDATE rides
	SET  start_otp = $1, ride_otp_generated_at = $2, arrived_at = $3, status = $4 WHERE id = $5 `

	_, err := database.Ridesafe.Exec(SQL, otp, time.Now(), arrivedAt, status, rideID)

	return err
}

func VerifyRideOTP(rideID string, otp string) error {

	SQL := `
	UPDATE rides
	SET otp_verified_at = NOW() WHERE id = $1 AND start_otp = $2 `

	result, err := database.Ridesafe.Exec(SQL, rideID, otp)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

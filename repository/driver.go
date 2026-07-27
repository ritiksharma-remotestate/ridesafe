package repository

import (
	"ridesafe/database"
	"ridesafe/models"

	"github.com/jmoiron/sqlx"
)

func CreateDriver(userID string) error {
	SQL := `INSERT INTO drivers (user_id) VALUES ($1)`

	_, err := database.Ridesafe.Exec(SQL, userID)
	if err != nil {
		return err
	}

	return nil
}

func GetDriverByUserID(userID string) (*models.Driver, error) {
	SQL := `select * from Drivers where user_id=$1`
	var getDriver models.Driver

	err := database.Ridesafe.Get(
		&getDriver, SQL, userID,
	)
	return &getDriver, err
}
func GetDriverLocation(userID string) (float64, float64, error) {
	SQL := ` SELECT current_latitude, current_longitude FROM drivers WHERE user_id = $1 `
	var lat, lon float64
	err := database.Ridesafe.QueryRow(SQL, userID).Scan(&lat, &lon)
	if err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}

func UpdateDriverLocation(userID string, latitude, longitude float64) error {
	SQL := `UPDATE Drivers SET current_latitude=$2 , current_longitude=$3 where user_id=$1 `
	_, err := database.Ridesafe.Exec(SQL, userID, latitude, longitude)
	return err

}

func SetDriverOnline(userID string, online bool) error {
	SQL := `UPDATE Drivers SET is_online=$2 where user_id=$1 `
	_, err := database.Ridesafe.Exec(SQL, userID, online)
	return err
}

func SetDriverAvailable(db sqlx.Ext, userID string, available bool) error {
	SQL := `UPDATE Drivers SET is_available=$2 where user_id=$1 `
	_, err := db.Exec(SQL, userID, available)
	return err
}

func GetAvailableDrivers(latitude, longitude float64) ([]models.Driver, error) {
	query := `
	SELECT *
	FROM (
		SELECT *,
			(
				6371 * acos(
					cos(radians($1))
					* cos(radians(current_latitude))
					* cos(radians(current_longitude) - radians($2))
					+ sin(radians($1))
					* sin(radians(current_latitude))
				)
			) AS distance
		FROM drivers
		WHERE is_online = TRUE
		  AND is_available = TRUE
		  AND current_latitude IS NOT NULL
		  AND current_longitude IS NOT NULL
	) d
	WHERE distance <= 4
	ORDER BY distance;`

	var drivers []models.Driver

	err := database.Ridesafe.Select(&drivers, query, latitude, longitude)
	if err != nil {
		return nil, err
	}

	return drivers, nil
}

package respository

import (
	"ridesafe/database"
	"ridesafe/models"
)

func CreateDriver(userID string) error {
	query := `
		INSERT INTO drivers (user_id)
		VALUES ($1)
	`

	_, err := database.Ridesafe.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}

func GetDriverByUserID(userID string) (*models.Driver,error){
	SQL:=`select * from Drivers where user_id=$1`
	var getDriver models.Driver

	err:= database.Ridesafe.Get(
		&getDriver,SQL,userID,
	)
	return &getDriver,err
}

func UpdateDriverLocation(userID string ,latitude ,longitude float64) error{
SQL:=`UPDATE Drivers SET current_latitude=$2 , current_longitude=$3 where user_id=$1 `
_, err := database.Ridesafe.Exec(SQL, userID,latitude,longitude)
	return err

}

func SetDriverOnline(userID string,online bool)error{
	SQL:=`UPDATE Drivers SET is_online=$2 where user_id=$1 `
_, err := database.Ridesafe.Exec(SQL, userID,online)
	return err
}

func SetDriverAvailable(userID string,available bool)error{
	SQL:=`UPDATE Drivers SET is_available=$2 where user_id=$1 `
_, err := database.Ridesafe.Exec(SQL, userID,available)
	return err
}

func GetAvailableDrivers()([]models.Driver,error){
SQL:=`SELECT * from Drivers where is_online= True and is_available=True`
var drivers []models.Driver
err := database.Ridesafe.Select(&drivers,SQL)
	return drivers,err
}



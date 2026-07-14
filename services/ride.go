package services 


import(
	"ridesafe/repository"
	"ridesafe/models"
	"time"
	"ridesafe/database"
	"errors"
	"ridesafe/utils"
	

)
func CreateRide(ride models.ride)(*models.Ride,error){
	user,err:=repository.GetUserByID(ride.PassengerID)
	if err!=nil{
		return nil, ErrUserNotFound
	}
	if models.Role(user.Role) != models.RoleUser {
		return nil,ErrNotAPassenger
	}
	//  can check if ther is any active ridess going on 

	distance := utils.CalculateDistance(
    ride.PickupLatitude,
    ride.PickupLongitude,
    ride.DestinationLatitude,
    ride.DestinationLongitude,
)

ride.Fare = utils.CalculateFare(distance)

ride.Status = models.RideRequested
	ride.RequestedAt = time.Now()

	newRide, err := repository.CreateRide(ride)
	if err != nil {
		return nil, err
	}

	return newRide, nil
}

func AcceptRide(rideID string,driverID string) error{

		ride,err:=repository.GetRideByID(rideID)
		if err!=nil{
			return ErrNotARide
		}
		if models.RideStatus(ride.Status)!=models.RideRequested{
			return ErrDriverAlreadyExists}

		driver,err:=repository.GetDriverByuserID(driverID)
		if err!=nil{
			return ErrNotADriver
		}

		if !driver.IsOnline{
			return ErrDriverAlreadyExists}
        if !driver.IsAvailable{
			return ErrDriverAlreadyExists}


		 return database.Tx(func(tx *sqlx.Tx) error {

        err := repository.AcceptRide(tx,
            rideID,
            driverID,
            time.Now(),
            models.RideAccepted,
        )
        if err != nil {
            return err
        }

        err = repository.SetDriverAvailable(
            tx,
            driverID,
            false,
        )
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
			return ErrRideNotFound
		}
		return err
	}

	if models.RideStatus(ride.Status) != models.RideAccepted {
		return ErrRideNotAccepted
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return ErrUnauthorizedDriver
	}


	err = repository.StartRide(
		rideID,
		time.Now(),
		models.RideStarted,
	)
	if err != nil {
		return err
	}

	return nil
	
}

func CompleteRide(rideID, driverID string) error {

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRideNotFound
		}
		return err
	}

	if models.RideStatus(ride.Status) != models.RideStarted {
		return ErrRideNotStarted
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return ErrUnauthorizedDriver
	}
	return database.Tx(func(tx *sqlx.Tx) error {

	err = repository.CompleteRide(
		tx,
		rideID,
		time.Now(),
		models.RideCompleted,
	)
	if err != nil {
		return err
	}

	err = repository.SetDriverAvailable(tx,driverID, true)
	if err != nil {
		return err
	}

	return nil
})
}
func CancelRide(rideID string) error {

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRideNotFound
		}
		return err
	}

	if models.RideStatus(ride.Status) == models.RideCompleted {
		return ErrRideAlreadyCompleted
	}

	if models.RideStatus(ride.Status) == models.RideCancelled {
		return ErrRideAlreadyCancelled
	}
	return database.Tx(func(tx *sqlx.Tx) error {
	err = repository.CancelRide(
		tx,
		rideID,
		time.Now(),
		models.RideCancelled,
	)
	if err != nil {
		return err
	}

	if ride.DriverID != nil {
		err = repository.SetDriverAvailable(tx,*ride.DriverID, true)
		if err != nil {
			return err
		}
	}

	return nil
})
}


func MarkRideArrived(rideID, driverID string) error {

	ride, err := repository.GetRideByID(rideID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRideNotFound
		}
		return err
	}

	if ride.Status != models.RideAccepted {
		return ErrRideNotAccepted
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return ErrUnauthorizedDriver
	}

	err = repository.MarkRideArrived(
		rideID,
		time.Now(),
		models.RideArrived,
	)
	if err != nil {
		return err
	}

	return nil
}


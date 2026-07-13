func CreateRide(ride models.ride)(*models.Ride,error){
	_,err:=respository.GetUserByID(ride.PassengerID)
	if err!=nil{
		return ErrUserNotFound
	}
	if user.Role != models.RolePassenger {
		return ErrNotAPassenger
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

func AcceptRide(rideID string,driverID string,acceptedAt time.Time,status models.RideStatus,
	) error{
		
	}




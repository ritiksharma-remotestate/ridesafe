func CreateRide(ride models.Ride)(*models.Ride,error){
	SQL:= `INSERT INTO rides ( passenger_id, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude, fare, status
)
VALUES ( $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;`

var newRide models.Ride

err := database.Ridesafe.Get( &newRide, SQL, ride.PassengerID, ride.PickupLatitude, ride.PickupLongitude, ride.DestinationLatitude, ride.DestinationLongitude, ride.Fare, ride.Status,
)
if err != nil {
    return nil, err
}
return &newRide,nil;
}



func AcceptRide(rideID string,driverID string,acceptedAt time.Time,status models.RideStatus,
	) error{
		 SQL:=`UPDATE rides SET driver_id = $1, accepted_at = $2, status = $3 WHERE id = $4;`
		
		
		
		_,err := database.Ridesafe.Exec(SQL, driverID, acceptedAt, status, rideID
)
if err != nil {
	return  err
}
return nil;
}




func StartRide( rideID string, started_at time.Time, status models.RideStatus,
	) error{
		SQL:=`UPDATE rides SET started_at= $1, status = $2 WHERE id = $3;`
		
		
		
		_,err := database.Ridesafe.Exec( SQL,   started_at,
    status,
    rideID,		)
		if err != nil {
			return  err
		}
		return nil;
	}
func CompleteRide(
		rideID string,
		
		completed_at time.Time,
		status models.RideStatus,
		) error{
			SQL:=`UPDATE rides
			SET
			
			completed_at= $1,
			status = $2
			WHERE id = $3;`
			
			
			
			_,err := database.Ridesafe.Exec(
				
				SQL,
				
				
				completed_at,
				status,
				rideID
			)
			if err != nil {
				return  err
			}
			return nil;
		}
		func CancelRide(
			rideID string,
			
			cancelled_at time.Time,
			status models.RideStatus,
			) error{
				SQL:=`UPDATE rides
				SET
				
				cancelled_at= $1,
				status = $2
				WHERE id = $3;`
				
				
				
				_,err := database.Ridesafe.Exec(
					
					SQL,
					
					
					cancelled_at,
					status,
					rideID
				)
				if err != nil {
					return  err
				}
				return nil;
			}
			
			func MarkRideArrived(
				rideID string,
				
				arrived_at time.Time,
				status models.RideStatus,
				) error{
					SQL:=`UPDATE rides
					SET
					
					arrived_at= $1,
					status = $2
					WHERE id = $3;`
					
					
					
					_,err := database.Ridesafe.Exec(
						
						SQL,
						
						arrived_at,
						status,
						rideID
					)
					if err != nil {
						return  err
					}
					return nil;
				}






				
				/*
				
				
				func AssignDriver(ride models.Ride)(*models.Ride,error){
					SQL:=`UPDATE Rides SET  (driver_id,fare,status) values $1,$2,$3 RETURNING *;`
				
					var newRide models.Ride
					err := database.Ridesafe.Exec(
					&newRide,
					SQL,
					ride.DriverID,
					
					ride.AcceptedAt,
					ride.Status,
				)
				if err != nil {
					return nil, err
				}
				return &newRide,nil;
				}
				func AssignDriver(ride models.Ride) error{
					SQL:=`UPDATE Rides Set driver_id=$1,accepted_at=$2 ,status=$3 where id=$4`
				
					
					_,err := database.Ridesafe.Exec(
					
					SQL,
					ride.DriverID,
					
					ride.AcceptedAt,
					ride.Status,
					ride.ID
				)
				if err != nil {
					return nil, err
				}
				return nil;
				}
				*/
				

func CreateDriver(userID string)error{

user,err:= repository.GetUserByID(userID)
if err!=nil{
	if errors.Is(err,sql.ErrNoRows){
		return ErrUserNotFound
	}
	return err
}
if user.Role != models.RoleDriver {
		return ErrNotADriver
	}

	driver, err := repository.GetDriverByUserID(userID)
	if err == nil && driver != nil {
		return ErrDriverAlreadyExists
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = repository.CreateDriver(userID)
	if err != nil {
		return err
	}

	return nil
}

func GetDriverByUserID(userID string)(*models.Driver,error){
	driver,err:=repository.GetDriverByuserID(userID)
	if err!=nil{
	if errors.Is(err,sql.ErrNoRows){
		return ErrUserNotFound
	}
	return nil,err

}
return driver,nil

}
func UpdateDriverLocation(userID string,latitude,longitude float64 ) error{
	_,err:=repository.GetDriverByuserID(userID)
	if err!=nil{
	if errors.Is(err,sql.ErrNoRows){
		return ErrUserNotFound
	}
	return err

}




	err=repository.UpdateDriverLocation(userID,latitude,longitude)
	if err!=nil{
		return err
	}
	return nil

}

func DriverOnline(userID string,online bool)error{
	_,err:=repository.GetDriverByUserID(userID)
	if err!=nil{
	if errors.Is(err,sql.ErrNoRows){
		return ErrUserNotFound
	}
	return err
	}
	err=repository.SetDriverOnline(userID,online)
	if err!=nil{
		return err
	}
	return nil


}
func DriverAvailable(userID string,available bool)error{
	_,err:=repository.GetDriverByUserID(userID)
	if err!=nil{
	if errors.Is(err,sql.ErrNoRows){
		return ErrUserNotFound
	}
	return err
}

	err=repository.SetDriverAvailable(userID,available)
	if err!=nil{
		return err
	}
	return nil


}

func GetAvailableDrivers() ([]models.Driver, error){
	
	drivers,err:=repository.GetAvailableDrivers()
	if err!=nil{
		return nil,err
	}
	return drivers,nil
}










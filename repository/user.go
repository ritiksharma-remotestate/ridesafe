package repository

import (
	"ridesafe/database"
	"ridesafe/models"
	"database/sql"
)

func RegisterUser(user models.User) (*models.User,error){
	 SQL:=`insert into users (name,email,hashed_password,role)
	 values ($1,$2,$3,$4) returning *`

	 var createdUser models.User

	 err:=database.Ridesafe.Get(
		&createdUser,SQL,user.Name,user.Email,user.HashedPassword,user.Role,
	 )
	 return &createdUser,err

}

func GetUserByID(userID string) (*models.User,error){
	SQL:=`Select  id, name,email,role from users where id=$1 and archived_at is NULL`
	var getUser models.User

	err:= database.Ridesafe.Get(
		&getUser,SQL,userID,
	)
	return &getUser,err
}

func GetUserByEmail(Email string) (*models.User,error){
	SQL:=`Select  id, name,email,role from users where email=$1 and archived_at is NULL`
	var getUser models.User

	err:= database.Ridesafe.Get(
		&getUser,SQL,Email,
	)
	return &getUser,err
}

func IsUserExists(email string) (bool, error) {
	query := `
		SELECT 1
		FROM users
		WHERE email=$1
		LIMIT 1
	`

	var result int

	err := database.Ridesafe.Get(
		&result,
		query,
		email,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}


func GetUserIDByPassword(email, password string) (string, error) {
	SQL := `SELECT
				u.id,
       			u.password
       		FROM
				users u
			WHERE
				u.archived_at IS NULL
				AND u.email = TRIM(LOWER($1))`
	var userID string
	var passwordHash string
	err := database.Ridesafe.QueryRowx(SQL, email).Scan(&userID, &passwordHash)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		return "", nil
	}
	// compare password
	if passwordErr := utils.CheckPassword(password, passwordHash); passwordErr != nil {
		return "", passwordErr
	}
	return userID, nil
}
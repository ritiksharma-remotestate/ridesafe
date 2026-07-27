package repository

import (
	"ridesafe/database"
	"ridesafe/models"
)

func RegisterUser(user models.RegisterRequest) (*models.User, error) {
	SQL := `insert into users (name,email,hashed_password,role)
	 values ($1,$2,$3,$4) RETURNING id, name, email, hashed_password, role, created_at, updated_at`

	var createdUser models.User

	err := database.Ridesafe.Get(
		&createdUser, SQL, user.Name, user.Email, user.Password, user.Role,
	)
	return &createdUser, err
}

func GetUserByID(userID string) (*models.User, error) {
	SQL := `Select  id, name,email,role from users where id=$1 and archived_at is NULL`
	var getUser models.User

	err := database.Ridesafe.Get(
		&getUser, SQL, userID,
	)
	return &getUser, err
}

func GetUserByEmail(Email string) (*models.User, error) {
	SQL := `Select  id, name,email,role , hashed_password from users where email=$1 and archived_at is NULL`
	var getUser models.User

	err := database.Ridesafe.Get(
		&getUser, SQL, Email,
	)
	return &getUser, err
}

func IsUserExists(email string) (bool, error) {
	SQL := `select count(*) > 0
       	from users where email=$1`

	var result bool
	err := database.Ridesafe.Get(&result, SQL, email)
	return result, err
}

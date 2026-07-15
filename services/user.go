package services

import (
	"database/sql"
	"errors"
	"ridesafe/error_custom"
	"golang.org/x/crypto/bcrypt"
	
	"ridesafe/models"
	"ridesafe/repository"
	"ridesafe/utils"
)

func RegisterUser(user models.RegisterRequest) (*models.User, error) {

	exists, err := repository.IsUserExists(user.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, error_custom.ErrUserAlreadyExists 
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword(
	// 	[]byte(user.Password),
	// 	bcrypt.DefaultCost,
	// )
hashedPassword,err:=utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	newUser, err := repository.RegisterUser(user)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func Login(email string, password string) (string, error) {

    user, err := repository.GetUserByEmail(email)

    if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return "", error_custom.ErrInvalidCredentials
    }
    return "", err
}


    err = bcrypt.CompareHashAndPassword(
        []byte(user.HashedPassword),
        []byte(password),
    )

    if err != nil {
        return "", error_custom.ErrInvalidCredentials
    }


    token, err := utils.JwtTokenCreate(user.ID,user.Name,email,user.Role)

    if err != nil {
        return "", err
    }


    return token, nil
}

func GetUserProfile(userID string) (*models.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error_custom.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}


// func RegisterUser(user models.User) (*models.User, error) {

// 	exists, err := repository.IsUserExists(user.Email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if exists {
// 		return nil, errors.New("email already exists")
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword(
// 		[]byte(user.Password),
// 		bcrypt.DefaultCost,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	user.HashedPassword = string(hashedPassword)

// 	newUser, err := repository.CreateUser(user)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return newUser, nil
// }

// func Login(email string, password string) (string, error) {

//     user, err := repository.GetUserByEmail(email)

//     if err != nil {
//         return "", err
//     }


//     err = bcrypt.CompareHashAndPassword(
//         []byte(user.HashedPassword),
//         []byte(password),
//     )

//     if err != nil {
//         return "", errors.New("invalid credentials")
//     }


//     token, err := GenerateJWT(user)

//     if err != nil {
//         return "", err
//     }


//     return token, nil
// }

// func GetUserProfile(userID string) (*models.User, error) {
// 	user, err := repository.GetUserByID(userID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, ErrUserNotFound
// 		}
// 		return nil, err
// 	}

// 	return user, nil
// }


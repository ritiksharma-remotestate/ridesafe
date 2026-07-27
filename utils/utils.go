package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"ridesafe/models"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"
)

var generator *shortid.Shortid

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

const generatorSeed = 1000

type Claims struct {
	UserID string      `json:"user_id"`
	Name   string      `json:"name"`
	Email  string      `json:"email"`
	Role   models.Role `json:"role"`

	jwt.RegisteredClaims
}

type FieldError struct {
	Err validator.ValidationErrors
}
type clientError struct {
	ID            string `json:"id"`
	MessageToUser string `json:"messageToUser"`
	DeveloperInfo string `json:"developerInfo"`
	Err           string `json:"error"`
	StatusCode    int    `json:"statusCode"`
	IsClientError bool   `json:"isClientError"`
}

func init() {
	n, err := rand.Int(rand.Reader, big.NewInt(generatorSeed))
	if err != nil {
		logrus.Panicf("failed to initialize utilities with random seed, %+v", err)
		return
	}

	g, err := shortid.New(1, shortid.DefaultABC, n.Uint64())

	if err != nil {
		logrus.Panicf("Failed to initialize utils package with error: %+v", err)
	}

	generator = g
}

func ParseBody(body io.Reader, out any) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

func EncodeJSONBody(resp http.ResponseWriter, data any) error {
	return json.NewEncoder(resp).Encode(data)
}

func RespondJSON(w http.ResponseWriter, statusCode int, body any) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			logrus.Errorf("failed to respond json with error:%v", err)
		}
	}
}

func newClientError(err error, statusCode int, messageToUser string, additionalInfoForDevs ...string) *clientError {
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	if additionalInfoJoined == "" {
		additionalInfoJoined = messageToUser
	}
	errorID, _ := generator.Generate()

	var errString string
	if err != nil {
		errString = err.Error()
	}
	return &clientError{
		ID:            errorID,
		MessageToUser: messageToUser,
		DeveloperInfo: additionalInfoJoined,
		Err:           errString,
		StatusCode:    statusCode,
		IsClientError: true,
	}
}

func RespondError(w http.ResponseWriter, statusCode int, err error, messageToUser string, additionalInfoForDevs ...string) {
	logrus.Errorf("status,%d,message:%s err:%v ", statusCode, messageToUser, err)
	clientError := newClientError(err, statusCode, messageToUser, additionalInfoForDevs...)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientError); err != nil {
		logrus.Errorf("failed to send error to caller with err:%v", err)
	}
}

func CheckValidation(i any) validator.ValidationErrors {
	v := validator.New()
	err := v.Struct(i)
	if err == nil {
		return nil
	}
	return err.(validator.ValidationErrors)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func JwtTokenCreate(userId, name, email string, role models.Role) (string, error) {
	claims := Claims{
		UserID: userId,
		Name:   name,
		Email:  email,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(

		tokenString,

		&Claims{},

		func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return jwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil

}

func CalculateFare(distance float64) float64 {
	fare := distance * 12
	return fare
}

func CalculateDistance(PickupLatitude, PickupLongitude, DestinationLatitude, DestinationLongitude float64) float64 {
	const EarthRadiusKm = 6371.0

	lat1Rad := PickupLatitude * math.Pi / 180
	lon1Rad := PickupLongitude * math.Pi / 180
	lat2Rad := DestinationLatitude * math.Pi / 180
	lon2Rad := DestinationLongitude * math.Pi / 180

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c
}
func GenerateOTP() (string, error) {
	maxNum := big.NewInt(9000)

	n, err := rand.Int(rand.Reader, maxNum)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%04d", n.Int64()+1000), nil
}

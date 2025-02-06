package helper

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/subhendu/go-auth/models"
)

func CheckUserType(r *http.Request, role string) error {
	userRole := r.Context().Value("role").(string)
	var err error = nil
	if userRole != role {
		err = errors.New("unauthorized, you are not able to access this data")
		return err
	}
	return err
}

func MatchUserTypeToUid(r *http.Request, user_id string) error {
	userId := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)
	var err error = nil
	if role == "USER" && userId != user_id {
		err = errors.New("unauthorized, you are not able to access this data")
		return err
	}
	err = CheckUserType(r, role)
	return err
}

type Claims struct {
	ID                   string `json:"id"`
	Email                string `json:"email"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Phone                string `json:"phone"`
	jwt.RegisteredClaims        // Embedding RegisteredClaims
}

var jwtSecret = []byte("12345")

func GenerateToken(user models.User, expiryTime time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiryTime)
	claims := &Claims{
		ID:        *user.User_type,
		Email:     *user.Emil,
		FirstName: *user.First_name,
		LastName:  *user.Last_name,
		Phone:     *user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

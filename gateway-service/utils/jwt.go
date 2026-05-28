package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	INVALID_AUTH_TOKEN = "invalid authorization token headers."
)

func ExtractToken(token string) (string, error) {
	reqToken := token
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New(INVALID_AUTH_TOKEN)
	}
	return splitToken[1], nil
}

func ExtractClaims(token string) (jwt.MapClaims, error) {
	tokenString, err := ExtractToken(token)
	if err != nil {
		return nil, err
	}

	hmacSecretString := os.Getenv("SERVICE_JWT_SECRET")
	hmacSecret := []byte(hmacSecretString)

	payload, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := payload.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func GetFullnameByToken(token string) (fullname string, success bool) {
	payload, check := ExtractClaims(token)
	if check != nil {
		return "", false
	}
	if payload["FirstName"] == nil {
		return "", false
	}
	fullname = payload["FirstName"].(string) + " " + payload["LastName"].(string)
	return fullname, true
}

func GetEmployeeID(token string) (id string, success bool) {
	payload, check := ExtractClaims(token)
	if check != nil {
		return "", false
	}
	if payload["ID"] == nil {
		return "", false
	}
	id = payload["ID"].(string)
	return id, true
}

func GetRoleID(token string) (id string, success bool) {
	payload, check := ExtractClaims(token)
	if check != nil {
		return "", false
	}
	if payload["RoleID"] == nil {
		return "", false
	}
	id = payload["RoleID"].(string)
	return id, true
}

func GetCompanyByToken(token string) (companyID string, success bool) {
	payload, check := ExtractClaims(token)
	if check != nil {
		return "", false
	}
	if payload["CompanyID"] == nil {
		return "", false
	}
	companyID = payload["CompanyID"].(string)
	return companyID, true
}

func GetSub(tokenString string) (sub string, success bool) {
	t := strings.Replace(tokenString, "Bearer ", "", 1)

	hmacSecretString := os.Getenv("SERVICE_JWT_SECRET")
	hmacSecret := []byte(hmacSecretString)

	val, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return sub, false
	}

	payload, ok := val.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	if payload["sub"] == nil {
		return "", false
	}
	return payload["sub"].(string), true
}

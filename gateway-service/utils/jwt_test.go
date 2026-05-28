package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(claims jwt.MapClaims) string {
	secret := "test-secret"
	os.Setenv("SERVICE_JWT_SECRET", secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestExtractToken(t *testing.T) {
	t.Run("ValidBearerToken", func(t *testing.T) {
		token, err := ExtractToken("Bearer abc123")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if token != "abc123" {
			t.Errorf("expected 'abc123', got '%s'", token)
		}
	})

	t.Run("MissingBearerPrefix", func(t *testing.T) {
		_, err := ExtractToken("abc123")
		if err == nil {
			t.Error("expected error for missing Bearer prefix")
		}
		if err.Error() != INVALID_AUTH_TOKEN {
			t.Errorf("expected '%s', got '%s'", INVALID_AUTH_TOKEN, err.Error())
		}
	})

	t.Run("EmptyString", func(t *testing.T) {
		_, err := ExtractToken("")
		if err == nil {
			t.Error("expected error for empty string")
		}
	})
}

func TestExtractClaims(t *testing.T) {
	t.Run("ValidToken", func(t *testing.T) {
		tokenString := generateTestToken(jwt.MapClaims{
			"sub":  "user-123",
			"exp":  time.Now().Add(time.Hour).Unix(),
			"name": "Test User",
		})

		claims, err := ExtractClaims("Bearer " + tokenString)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if claims["sub"] != "user-123" {
			t.Errorf("expected sub 'user-123', got '%v'", claims["sub"])
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		_, err := ExtractClaims("Bearer invalid-token")
		if err == nil {
			t.Error("expected error for invalid token")
		}
	})

	t.Run("MissingBearer", func(t *testing.T) {
		_, err := ExtractClaims("no-bearer-prefix")
		if err == nil {
			t.Error("expected error for missing Bearer")
		}
	})
}

func TestGetSub(t *testing.T) {
	t.Run("ValidTokenWithSub", func(t *testing.T) {
		tokenString := generateTestToken(jwt.MapClaims{
			"sub": "user-456",
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		sub, ok := GetSub("Bearer " + tokenString)
		if !ok {
			t.Error("expected success")
		}
		if sub != "user-456" {
			t.Errorf("expected 'user-456', got '%s'", sub)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		_, ok := GetSub("Bearer invalid")
		if ok {
			t.Error("expected failure for invalid token")
		}
	})

	t.Run("TokenWithoutSub", func(t *testing.T) {
		tokenString := generateTestToken(jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		_, ok := GetSub("Bearer " + tokenString)
		if ok {
			t.Error("expected failure for missing sub claim")
		}
	})
}

func TestGetFullnameByToken(t *testing.T) {
	t.Run("ValidTokenWithName", func(t *testing.T) {
		tokenString := generateTestToken(jwt.MapClaims{
			"FirstName": "John",
			"LastName":  "Doe",
			"exp":       time.Now().Add(time.Hour).Unix(),
		})

		fullname, ok := GetFullnameByToken("Bearer " + tokenString)
		if !ok {
			t.Error("expected success")
		}
		if fullname != "John Doe" {
			t.Errorf("expected 'John Doe', got '%s'", fullname)
		}
	})

	t.Run("MissingFirstName", func(t *testing.T) {
		tokenString := generateTestToken(jwt.MapClaims{
			"LastName": "Doe",
			"exp":      time.Now().Add(time.Hour).Unix(),
		})

		_, ok := GetFullnameByToken("Bearer " + tokenString)
		if ok {
			t.Error("expected failure for missing FirstName")
		}
	})
}

func TestGetEmployeeID(t *testing.T) {
	t.Run("ValidToken", func(t *testing.T) {
		tokenString := generateTestToken(jwt.MapClaims{
			"ID":  "emp-789",
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		id, ok := GetEmployeeID("Bearer " + tokenString)
		if !ok {
			t.Error("expected success")
		}
		if id != "emp-789" {
			t.Errorf("expected 'emp-789', got '%s'", id)
		}
	})

	t.Run("MissingID", func(t *testing.T) {
		tokenString := generateTestToken(jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		_, ok := GetEmployeeID("Bearer " + tokenString)
		if ok {
			t.Error("expected failure for missing ID")
		}
	})
}

func TestGetCompanyByToken(t *testing.T) {
	t.Run("ValidToken", func(t *testing.T) {
		tokenString := generateTestToken(jwt.MapClaims{
			"CompanyID": "comp-001",
			"exp":       time.Now().Add(time.Hour).Unix(),
		})

		companyID, ok := GetCompanyByToken("Bearer " + tokenString)
		if !ok {
			t.Error("expected success")
		}
		if companyID != "comp-001" {
			t.Errorf("expected 'comp-001', got '%s'", companyID)
		}
	})
}

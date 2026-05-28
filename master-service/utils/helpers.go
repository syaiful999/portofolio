package utils

import (
	"context"
	"database/sql"
	"encoding/base64"
	convert "encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	yaml "github.com/asim/go-micro/plugins/config/encoder/yaml/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source/file"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	secretKey     = os.Getenv("SERVICE_JWT_SECRET")
)

type TokenValue struct {
	ID     uuid.UUID `json:"id" db:"id"`
	RoleID uuid.UUID `json:"role_id" db:"role_id"`
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetKey(keyName string) (string, error) {
	if err := godotenv.Load(); err != nil {
		PushLogf("", "error load .env", "")
		return "", err
	}
	return os.Getenv(keyName), nil

}

func FindInArray(needle interface{}, hystack interface{}) (int, bool) {
	switch key := needle.(type) {
	case string:
		for i, item := range hystack.([]string) {
			if key == item {
				return i, true
			}
		}
	case int:
		for i, item := range hystack.([]int) {
			if key == item {
				return i, true
			}
		}
	case int64:
		for i, item := range hystack.([]int64) {
			if key == item {
				return i, true
			}
		}
	default:
		return -1, false
	}
	return -1, false
}

func ReadJsonFile(file *multipart.FileHeader) (string, error) {
	var result string
	read, errOpen := file.Open()
	if errOpen != nil {
		return result, nil
	}
	defer read.Close()

	buf := make([]byte, file.Size)

	len, _ := read.Read(buf)
	b := string(buf)
	if errOpen != nil {
		return result, errOpen
	}

	if len == 0 {
		return result, errors.New("invalid json file")
	}

	return b, nil
}

func ReadXMlFile(file *multipart.FileHeader) (string, error) {
	var result string
	read, errOpen := file.Open()
	if errOpen != nil {
		return result, nil
	}
	defer read.Close()

	buf := make([]byte, file.Size)

	len, _ := read.Read(buf)
	b := string(buf)
	if errOpen != nil {
		return result, errOpen
	}

	if len == 0 {
		return result, errors.New("invalid xml file")
	}

	return b, nil
}

func TrimLastChar(data string) string {
	if len(data) == 0 {
		return ""
	}
	data = data[0 : len(data)-1]
	return data
}

func ArrToStrDelimiter(arrs []string, delimiter string) string {
	arrCheck := []string{}
	result := ""
	if delimiter == "" {
		delimiter = ","
	}
	for _, arr := range arrs {
		if _, err := FindInArray(arrCheck, arr); !err {
			arrCheck = append(arrCheck, arr)
			result += "'" + arr + "'" + delimiter
		}
	}
	return strings.TrimSuffix(result, delimiter)
}

func ConvStringToInt(data string) int {
	a, b := strconv.Atoi(data)
	if b != nil {
		return 0
	}
	return a
}

func ConvStringToFloat64(data string) float64 {
	a, b := strconv.ParseFloat(data, 64)
	if b != nil {
		return 0
	}
	return a
}

func CheckRubbishString(data string, rubbish string) string {
	pText := strings.Trim(data, rubbish)
	pTextCheck := strings.Trim(pText, " ")
	if len(pTextCheck) == 0 {
		return ""
	}
	return pText
}

func Percentage(a int, b int) float64 {
	if b == 0 {
		b = 1
	}
	c := fmt.Sprintf("%.2f", float64(a)/float64(b)*100)
	d, _ := strconv.ParseFloat(c, 64)
	return d
}

func ReadConfigYaml(key string) (interface{}, error) {
	enc := yaml.NewEncoder()

	// new config
	c, err := config.NewConfig(
		config.WithReader(
			json.NewReader(
				reader.WithEncoder(enc),
			),
		),
	)
	if err != nil {
		return "", err
	}

	if err := c.Load(
		file.NewSource(
			file.WithPath("./config.yaml"),
		),
	); err != nil {
		return "", err
	}

	var data map[string]interface{}

	if err := c.Get("hosts", "database").Scan(&data); err != nil {
		return "", err
	}

	for index, v := range data {
		if index == key {
			return v, nil
		}
	}

	return "", errors.New("key not found")
}

func GetFullNameFromMetaData(ctx context.Context) string {
	//Get User Fullname From Metadata
	Fullname := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ufullName, ok := md["x-user-fullname"]; ok {
			Fullname = strings.Join(ufullName, ",")
		}
	}
	return Fullname
}

func GetCompanyIDFromMetaData(ctx context.Context) string {
	//Get User Fullname From Metadata
	companyId := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ucompanyId, ok := md["x-company-id"]; ok {
			companyId = strings.Join(ucompanyId, ",")
		}
	}
	return companyId
}

func ConvertSqlcToProto(from interface{}, to interface{}) error {

	//Get Marshal
	data, err := convert.Marshal(from)
	if err != nil {
		return err
	}

	//Unmarshal data
	err = convert.Unmarshal(data, to)
	if err != nil {
		return err
	}

	return err
}

func HandleMandatoryField(value string) bool {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ") == ""
}

func UploadToStatic(base64Decoded string, directory string, imageName string) (string, error) {

	decodedBytes, err := base64.StdEncoding.DecodeString(base64Decoded)
	if err != nil {
		fmt.Println("Error decoding base64:", err)

		return "", err
	}
	// Specify the directory to save the decoded content
	outputDirectory := fmt.Sprintf("./static/%s", directory)
	additionalName := removeCharacters(ConvTimeToString(time.Now()), "-", " ", ":") + "."

	imageName = strings.ReplaceAll(imageName, ".", additionalName)
	outputFilePath := outputDirectory + "/" + imageName

	// Create the output directory if it doesn't exist
	existDirectory, err := directoryExists(outputFilePath)
	if err != nil {
		fmt.Println("error checking directory", err)
	}
	if !existDirectory {
		if err := os.MkdirAll(outputDirectory, 0755); err != nil {
			fmt.Println("Error creating directory:", err)

			return "", err
		}
	}

	// Write the decoded content to a file
	err = os.WriteFile(outputFilePath, decodedBytes, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)

		return "", err
	}

	fmt.Printf("Decoded content saved to: %s\n", outputFilePath)
	// */
	return outputFilePath, nil

}

func ConvStringToNullString(value string) sql.NullString {
	var result sql.NullString
	if value == "" {
		return result
	} else {
		result = sql.NullString{String: value, Valid: true}
		return result
	}
}

func ConvInt32toNullInt32(value int32) sql.NullInt32 { return sql.NullInt32{Int32: value, Valid: true} }

func ConvStringToTime(dateString string) (time.Time, error) {
	// Define the layout of the input string
	layout := "02-01-2006"

	// Parse the string into a time.Time object
	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("[ConvStringToTime] Error parsing time:", dateString, err)
		var nullTime time.Time
		return nullTime, err
	}
	return parsedTime, nil
}

func ConvStringToNullTime(dateString string) (sql.NullTime, error) {
	// Define the layout of the input string
	var result sql.NullTime
	if dateString == "" {
		return result, nil
	}

	layout := "02-01-2006"
	// Parse the string into a time.Time object
	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("[ConvStringToNullTime] Error parsing time:", dateString, err)
		return result, err
	}
	result.Time = parsedTime
	result.Valid = true
	return result, nil
}

func ConvTimeStampToNullTime(dateTime *timestamppb.Timestamp) (sql.NullTime, error) {
	var result sql.NullTime
	if dateTime == nil {
		return result, nil
	}

	result.Time = dateTime.AsTime()
	result.Valid = true
	return result, nil
}

// converts time.Time to string with format "DD-MM-YYYY HH:MM:SS"
func ConvTimeToString(dateTime time.Time) string {

	if dateTime.IsZero() {
		return ""
	}
	// Define the layout for the desired string representation
	layout := "02-01-2006 15:04:05"

	// Format the time.Time object to a string using the defined layout
	timeString := dateTime.Format(layout)

	// Print the formatted string
	return timeString
}

// converts time.Time to string with format "DD-MM-YYYY"
func ConvTimeToStringDate(dateTime time.Time) string {

	if dateTime.IsZero() {
		return ""
	}
	// Define the layout for the desired string representation
	layout := "02-01-2006"

	// Format the time.Time object to a string using the defined layout
	timeString := dateTime.Format(layout)

	// Print the formatted string
	return timeString
}

func removeCharacters(text string, charToremove ...string) string {
	for _, character := range charToremove {
		text = strings.ReplaceAll(text, character, "")
	}
	return text

}

func directoryExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Directory does not exist
		return false, nil
	} else if err != nil {
		// An error occurred while checking
		return false, err
	}

	// Directory exists
	return true, nil
}

// hash-ing password from client to be encrypted using GenerateFromPassword from crypto/bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// comparing password & hashed password from client to be encrypted using GenerateFromPassword from crypto/bcrypt
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("error when get metadata")
	}
	arrToken := md["x-user-token"]
	if len(arrToken) == 0 || arrToken[0] == "" {
		return "", errors.New("unauthorized access")
	}

	return arrToken[0], nil
}

func TokenDecoder(tokenString string) (*TokenValue, error) {
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !tkn.Valid {
		return nil, errors.New("invalid token")
	}
	id, _ := claims.GetSubject()
	_ = id
	tokenObj := &TokenValue{
		ID:     uuid.MustParse(claims["id"].(string)),
		RoleID: uuid.MustParse(claims["role_id"].(string)),
	}
	return tokenObj, nil
}

func HandleToken(ctx context.Context) (*TokenValue, error) {
	var err error
	token, err := GetUserToken(ctx)
	if err != nil {
		return nil, err
	}
	tokenValue, err := TokenDecoder(token)
	if err != nil {
		return nil, err
	}
	return tokenValue, nil
}

func ParseUUID(id string) uuid.UUID {
	parsedID, _ := uuid.Parse(id)
	return parsedID
}

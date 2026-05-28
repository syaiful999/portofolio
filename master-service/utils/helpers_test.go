package utils

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestToSnakeCase(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"CamelCase", "MyVariableName", "my_variable_name"},
		{"Already snake_case", "already_snake", "already_snake"},
		{"Single word", "word", "word"},
		{"With numbers", "Var1Name2", "var1_name2"},
		{"Empty string", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ToSnakeCase(tc.input))
		})
	}
}

func TestFindInArray(t *testing.T) {
	t.Run("String slice", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		idx, found := FindInArray("banana", slice)
		assert.True(t, found)
		assert.Equal(t, 1, idx)

		_, found = FindInArray("grape", slice)
		assert.False(t, found)
	})

	t.Run("Int slice", func(t *testing.T) {
		slice := []int{10, 20, 30}
		idx, found := FindInArray(30, slice)
		assert.True(t, found)
		assert.Equal(t, 2, idx)

		_, found = FindInArray(40, slice)
		assert.False(t, found)
	})

	t.Run("Int64 slice", func(t *testing.T) {
		slice := []int64{100, 200, 300}
		idx, found := FindInArray(int64(100), slice)
		assert.True(t, found)
		assert.Equal(t, 0, idx)

		_, found = FindInArray(int64(400), slice)
		assert.False(t, found)
	})

	t.Run("Unsupported type", func(t *testing.T) {
		slice := []float64{1.1, 2.2}
		idx, found := FindInArray(1.1, slice)
		assert.False(t, found)
		assert.Equal(t, -1, idx)
	})
}

func TestTrimLastChar(t *testing.T) {
	assert.Equal(t, "hello", TrimLastChar("hello,"))
	assert.Equal(t, "", TrimLastChar("a"))
	assert.Equal(t, "", TrimLastChar(""))
}

func TestArrToStrDelimiter(t *testing.T) {
	arr := []string{"a", "b", "c"}
	expected := "'a','b','c'"
	assert.Equal(t, expected, ArrToStrDelimiter(arr, ","))

	expectedWithSemicolon := "'a';'b';'c'"
	assert.Equal(t, expectedWithSemicolon, ArrToStrDelimiter(arr, ";"))

	// Test with default delimiter
	assert.Equal(t, expected, ArrToStrDelimiter(arr, ""))
}

func TestConvStringToInt(t *testing.T) {
	assert.Equal(t, 123, ConvStringToInt("123"))
	assert.Equal(t, 0, ConvStringToInt("abc"))
	assert.Equal(t, 0, ConvStringToInt(""))
}

func TestConvStringToFloat64(t *testing.T) {
	assert.Equal(t, 123.45, ConvStringToFloat64("123.45"))
	assert.Equal(t, 0.0, ConvStringToFloat64("abc"))
	assert.Equal(t, 0.0, ConvStringToFloat64(""))
}

func TestCheckRubbishString(t *testing.T) {
	assert.Equal(t, "hello", CheckRubbishString("...hello...", "."))
	assert.Equal(t, "world", CheckRubbishString(" world ", " "))
	assert.Equal(t, "", CheckRubbishString("---", "-"))
}

func TestPercentage(t *testing.T) {
	assert.Equal(t, 50.0, Percentage(1, 2))
	assert.Equal(t, 33.33, Percentage(1, 3))
	assert.Equal(t, 0.0, Percentage(0, 100))
	// Test division by zero
	assert.Equal(t, 100.0, Percentage(1, 0))
}

func TestGetFromMetaData(t *testing.T) {
	ctx := context.Background()
	md := metadata.New(map[string]string{
		"x-user-fullname": "John Doe",
		"x-company-id":    "company-123",
	})
	ctxWithMD := metadata.NewIncomingContext(ctx, md)

	t.Run("Get FullName", func(t *testing.T) {
		fullName := GetFullNameFromMetaData(ctxWithMD)
		assert.Equal(t, "John Doe", fullName)
	})

	t.Run("Get CompanyID", func(t *testing.T) {
		companyID := GetCompanyIDFromMetaData(ctxWithMD)
		assert.Equal(t, "company-123", companyID)
	})

	t.Run("Metadata not present", func(t *testing.T) {
		fullName := GetFullNameFromMetaData(context.Background())
		assert.Equal(t, "", fullName)
	})
}

func TestHandleMandatoryField(t *testing.T) {
	assert.True(t, HandleMandatoryField(""))
	assert.True(t, HandleMandatoryField("   "))
	assert.False(t, HandleMandatoryField("hello"))
	assert.False(t, HandleMandatoryField("  hello world  "))
}

func TestConvStringToNullString(t *testing.T) {
	ns := ConvStringToNullString("hello")
	assert.True(t, ns.Valid)
	assert.Equal(t, "hello", ns.String)

	nsEmpty := ConvStringToNullString("")
	assert.False(t, nsEmpty.Valid)
}

func TestConvStringToTime(t *testing.T) {
	layout := "02-01-2006"
	expectedTime, _ := time.Parse(layout, "17-12-2025")

	parsedTime, err := ConvStringToTime("17-12-2025")
	assert.NoError(t, err)
	assert.Equal(t, expectedTime, parsedTime)

	_, err = ConvStringToTime("invalid-date")
	assert.Error(t, err)
}

func TestConvStringToNullTime(t *testing.T) {
	nt, err := ConvStringToNullTime("17-12-2025")
	assert.NoError(t, err)
	assert.True(t, nt.Valid)

	ntEmpty, err := ConvStringToNullTime("")
	assert.NoError(t, err)
	assert.False(t, ntEmpty.Valid)
}

func TestConvTimeToString(t *testing.T) {
	tm, _ := time.Parse("2006-01-02 15:04:05", "2025-12-17 10:30:00")
	assert.Equal(t, "17-12-2025 10:30:00", ConvTimeToString(tm))
	assert.Equal(t, "", ConvTimeToString(time.Time{}))
}

func TestConvTimeToStringDate(t *testing.T) {
	tm, _ := time.Parse("2006-01-02 15:04:05", "2025-12-17 10:30:00")
	assert.Equal(t, "17-12-2025", ConvTimeToStringDate(tm))
	assert.Equal(t, "", ConvTimeToStringDate(time.Time{}))
}

func TestPasswordHashing(t *testing.T) {
	password := "my-secret-password"
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	assert.True(t, CheckPasswordHash(password, hash))
	assert.False(t, CheckPasswordHash("wrong-password", hash))
}

func TestParseUUID(t *testing.T) {
	id := uuid.New()
	parsed := ParseUUID(id.String())
	assert.Equal(t, id, parsed)

	invalidParsed := ParseUUID("not-a-uuid")
	assert.Equal(t, uuid.Nil, invalidParsed)
}

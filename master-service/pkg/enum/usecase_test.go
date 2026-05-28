package enum

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	pb "moyo-master-service/pkg/enum/proto"
	"moyo-master-service/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockEnumRepository is a mock implementation of IEnumRepository for testing.
type MockEnumRepository struct {
	MockGetEnum      func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEnum, int32, error)
	MockGetEnumById  func(ctx context.Context, id uuid.UUID) (MasterEnum, error)
	MockUpdateEnum   func(ctx context.Context, value MasterEnum) (MasterEnum, error)
	MockGetEnumExcel func() ([]MasterEnum, error)
}

func (m *MockEnumRepository) GetEnum(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEnum, int32, error) {
	if m.MockGetEnum != nil {
		return m.MockGetEnum(ctx, skip, take, filter, sort)
	}
	return nil, 0, nil
}

func (m *MockEnumRepository) GetEnumById(ctx context.Context, id uuid.UUID) (MasterEnum, error) {
	if m.MockGetEnumById != nil {
		return m.MockGetEnumById(ctx, id)
	}
	return MasterEnum{}, nil
}

func (m *MockEnumRepository) UpdateEnum(ctx context.Context, value MasterEnum) (MasterEnum, error) {
	if m.MockUpdateEnum != nil {
		return m.MockUpdateEnum(ctx, value)
	}
	return MasterEnum{}, nil
}

func (m *MockEnumRepository) GetEnumExcel() ([]MasterEnum, error) {
	if m.MockGetEnumExcel != nil {
		return m.MockGetEnumExcel()
	}
	return nil, nil
}

func TestGetEnums(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)

	// Mock implementation for GetEnum
	mockRepo.MockGetEnum = func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEnum, int32, error) {
		return []MasterEnum{{ID: uuid.New(), IsActive: true, CreatedBy: "enum", EnumValue: "Enum1"}}, 1, nil
	}

	req := &pb.GetEnumRequest{
		Skip:   0,
		Take:   1,
		Filter: "",
		Sort:   "",
	}

	res := &pb.GetEnumResponse{}
	err := useCase.GetEnums(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, int32(1), res.CountData)
}

func TestGetEnums_Error(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)

	// Mock implementation for GetEnum with error
	mockRepo.MockGetEnum = func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEnum, int32, error) {
		return nil, 0, errors.New("error fetching enums")
	}

	req := &pb.GetEnumRequest{
		Skip:   0,
		Take:   1,
		Filter: "",
		Sort:   "",
	}

	res := &pb.GetEnumResponse{}
	err := useCase.GetEnums(context.Background(), req, res)

	// usecase returns nil but sets response error fields
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusInternalServerError), res.ErrorCode)
	assert.Nil(t, res.Data)
	assert.Equal(t, int32(0), res.CountData)
}

func TestUpdateEnum(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)

	expectedID := uuid.New()

	mockRepo.MockUpdateEnum = func(ctx context.Context, value MasterEnum) (MasterEnum, error) {
		return MasterEnum{ID: expectedID, IsActive: value.IsActive, EnumValue: "UpdatedEnum_Updated"}, nil
	}

	req := &pb.UpdateEnumRequest{
		Id:        expectedID.String(),
		EnumValue: "UpdatedEnum",
	}

	res := &pb.UpdateEnumResponse{}
	token := &utils.TokenValue{ID: expectedID}
	err := useCase.UpdateEnum(context.Background(), req, res, token)

	assert.NoError(t, err)
	assert.Equal(t, expectedID.String(), res.Data.Id)
}

func TestUpdateEnum_InvalidID(t *testing.T) {
	// Test repository error during update
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)

	expectedID := uuid.New()
	req := &pb.UpdateEnumRequest{
		Id:        expectedID.String(),
		EnumValue: "SomeValue",
		EnumType:  "some_type",
	}
	res := &pb.UpdateEnumResponse{}
	token := &utils.TokenValue{ID: expectedID}

	mockRepo.MockUpdateEnum = func(ctx context.Context, value MasterEnum) (MasterEnum, error) {
		return MasterEnum{}, errors.New("db error")
	}

	err := useCase.UpdateEnum(context.Background(), req, res, token)
	// usecase returns nil but sets response error fields
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusInternalServerError), res.ErrorCode)
}

func TestUpdateEnum_MissingMandatoryField(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)

	expectedID := uuid.New()
	req := &pb.UpdateEnumRequest{
		Id: expectedID.String(),
	}
	res := &pb.UpdateEnumResponse{}
	token := &utils.TokenValue{}

	useCase.UpdateEnum(context.Background(), req, res, token)

	err := errors.New(res.ErrorMessage)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "enum value is required")
}

func TestGetEnumById(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)
	mockEnumID := uuid.New()

	expectedEnum := MasterEnum{
		ID:        mockEnumID,
		EnumValue: "Enum1",
	}

	mockRepo.MockGetEnumById = func(ctx context.Context, id uuid.UUID) (MasterEnum, error) {
		return expectedEnum, nil
	}

	req := &pb.GetEnumByIDRequest{
		Id: mockEnumID.String(),
	}

	res := &pb.GetEnumByIDResponse{}
	err := useCase.GetEnumById(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, expectedEnum.ID.String(), req.Id)
	assert.NotNil(t, res.Data)
	assert.Equal(t, expectedEnum.EnumValue, res.Data.EnumValue)
}

func TestGetEnumById_ErrorCases(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)
	mockEnumID := uuid.New()

	mockRepo.MockGetEnumById = func(ctx context.Context, id uuid.UUID) (MasterEnum, error) {
		return MasterEnum{}, fmt.Errorf("enum id %s not found", mockEnumID)
	}

	req := &pb.GetEnumByIDRequest{
		Id: mockEnumID.String(),
	}

	res := &pb.GetEnumByIDResponse{}
	err := useCase.GetEnumById(context.Background(), req, res)

	// usecase returns nil but populates response with error info
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	// repository returned a not-found like error; ensure response indicates failure
	assert.NotNil(t, res.ErrorMessage)
	assert.Nil(t, res.Data)
}

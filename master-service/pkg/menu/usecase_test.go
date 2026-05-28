package menu

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	pb "moyo-master-service/pkg/menu/proto"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockMenuRepository is a mock implementation of IMenuRepository for testing.
type MockMenuRepository struct {
	MockGetMenu               func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterMenu, int32, error)
	MockGetMenuById           func(ctx context.Context, id uuid.UUID) ([]MasterMenu, error)
	MockCreateMenu            func(ctx context.Context, value MasterMenu) (MasterMenu, error)
	MockUpdateMenu            func(ctx context.Context, value MasterMenu) (MasterMenu, error)
	MockDeleteMenu            func(ctx context.Context, id, modifiedBy uuid.UUID) error
	MockGetMasterMenuByName   func(ctx context.Context, menu_name string) (int32, error)
	MockGetMasterMenuValidate func(ctx context.Context, menu_name string, id uuid.UUID) (int32, error)
}

func (m *MockMenuRepository) GetMenu(ctx context.Context, skip, take int32, filter, sort string) ([]MasterMenu, int32, error) {
	if m.MockGetMenu != nil {
		return m.MockGetMenu(ctx, skip, take, filter, sort)
	}
	return nil, 0, nil
}

func (m *MockMenuRepository) GetMenuById(ctx context.Context, id uuid.UUID) ([]MasterMenu, error) {
	if m.MockGetMenuById != nil {
		return m.MockGetMenuById(ctx, id)
	}
	return nil, nil
}

func (m *MockMenuRepository) CreateMenu(ctx context.Context, value MasterMenu) (MasterMenu, error) {
	if m.MockCreateMenu != nil {
		return m.MockCreateMenu(ctx, value)
	}
	return MasterMenu{}, nil
}

func (m *MockMenuRepository) UpdateMenu(ctx context.Context, value MasterMenu) (MasterMenu, error) {
	if m.MockUpdateMenu != nil {
		return m.MockUpdateMenu(ctx, value)
	}
	return MasterMenu{}, nil
}

func (m *MockMenuRepository) DeleteMenu(ctx context.Context, id, modifiedBy uuid.UUID) error {
	if m.MockDeleteMenu != nil {
		return m.MockDeleteMenu(ctx, id, modifiedBy)
	}
	return nil
}

func (m *MockMenuRepository) GetMasterMenuByName(ctx context.Context, menu_name string) (int32, error) {
	if m.MockGetMasterMenuByName != nil {
		return m.MockGetMasterMenuByName(ctx, menu_name)
	}
	return 0, nil
}

func (m *MockMenuRepository) GetMasterMenuValidate(ctx context.Context, menu_name string, id uuid.UUID) (int32, error) {
	if m.MockGetMasterMenuValidate != nil {
		return m.MockGetMasterMenuValidate(ctx, menu_name, id)
	}
	return 0, nil
}

func TestGetMenus(t *testing.T) {
	mockRepo := new(MockMenuRepository)
	useCase := NewUseCaseMenu(mockRepo)

	// Mock implementation for GetMenu
	mockRepo.MockGetMenu = func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterMenu, int32, error) {
		return []MasterMenu{{ID: uuid.New(), IsActive: true, CreatedBy: "menu", MenuCode: "Menu1"}}, 1, nil
	}

	req := &pb.GetMenuRequest{
		Skip:   0,
		Take:   1,
		Filter: "",
		Sort:   "",
	}

	res := &pb.GetMenuResponse{}
	err := useCase.GetMenus(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, int32(1), res.CountData)
}

func TestGetMenus_InvalidTake(t *testing.T) {
	mockRepo := new(MockMenuRepository)
	useCase := NewUseCaseMenu(mockRepo)

	req := &pb.GetMenuRequest{
		Skip:   0,
		Take:   -1,
		Filter: "",
		Sort:   "",
	}

	res := &pb.GetMenuResponse{}
	useCase.GetMenus(context.Background(), req, res)
	err := errors.New(res.ErrorMessage)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parameter take, Req.Take = -1")
	assert.Nil(t, res.Data)
}

func TestGetMenus_ErrorFromRepository(t *testing.T) {

	mockRepo := new(MockMenuRepository)
	useCase := NewUseCaseMenu(mockRepo)

	mockRepo.MockGetMenu = func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterMenu, int32, error) {
		return []MasterMenu{}, 0, sql.ErrNoRows
	}
	req := &pb.GetMenuRequest{
		Skip:   0,
		Take:   1,
		Filter: "",
		Sort:   "",
	}

	res := &pb.GetMenuResponse{}

	useCase.GetMenus(context.Background(), req, res)

	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusInternalServerError), res.ErrorCode)
	assert.Contains(t, res.ErrorMessage, "failed to GetMenus:")
	assert.Nil(t, res.Data)
}

func TestGetMenus_ErrorGeneratingConditions(t *testing.T) {
	mockRepo := new(MockMenuRepository)
	useCase := NewUseCaseMenu(mockRepo)

	req := &pb.GetMenuRequest{
		Sort:   "invalid-sort",
		Filter: "invalid-filter",
	}

	// Mock implementation for GetMenu
	mockRepo.MockGetMenu = func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterMenu, int32, error) {
		return []MasterMenu{}, 0, errors.New("Expected an error for generating conditions, but got nil.")
	}

	res := &pb.GetMenuResponse{}

	useCase.GetMenus(context.Background(), req, res)

	err := errors.New(res.ErrorMessage)

	assert.Error(t, err, "Expected an error for generating conditions, but got nil.")
	assert.Nil(t, res.Data, "Data should be nil when there's an error.")
	assert.Empty(t, res.CountData, "CountData should be empty when there's an error.")

}

func TestGetMenuById(t *testing.T) {
	mockRepo := new(MockMenuRepository)
	useCase := NewUseCaseMenu(mockRepo)
	mockMenuID := uuid.New()

	expectedMenu := MasterMenu{
		ID:       mockMenuID,
		MenuCode: "Menu1",
	}

	mockRepo.MockGetMenuById = func(ctx context.Context, id uuid.UUID) ([]MasterMenu, error) {
		return []MasterMenu{{ID: uuid.New(), IsActive: true, CreatedBy: "menu", MenuCode: "Menu1"}}, nil
	}

	req := &pb.GetMenuByIDRequest{
		Id: mockMenuID.String(),
	}

	res := &pb.GetMenuByIDResponse{}
	err := useCase.GetMenuById(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, expectedMenu.ID.String(), req.Id)
	assert.Equal(t, expectedMenu.MenuCode, res.Data.MenuCode)
}

func TestGetMenuById_ErrorCases(t *testing.T) {
	mockRepo := new(MockMenuRepository)
	useCase := NewUseCaseMenu(mockRepo)
	mockMenuID := uuid.New()

	mockRepo.MockGetMenuById = func(ctx context.Context, id uuid.UUID) ([]MasterMenu, error) {
		return nil, fmt.Errorf("menu id %s not found", mockMenuID)
	}

	req := &pb.GetMenuByIDRequest{
		Id: mockMenuID.String(),
	}

	res := &pb.GetMenuByIDResponse{}
	useCase.GetMenuById(context.Background(), req, res)

	err := errors.New(res.ErrorMessage)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), err.Error())
	assert.Nil(t, res.Data)
}

func TestGetMenuById_NotFound(t *testing.T) {
	mockRepo := new(MockMenuRepository)
	useCase := NewUseCaseMenu(mockRepo)
	mockMenuID := uuid.New()

	mockRepo.MockGetMenuById = func(ctx context.Context, id uuid.UUID) ([]MasterMenu, error) {
		return nil, sql.ErrNoRows
	}

	req := &pb.GetMenuByIDRequest{
		Id: mockMenuID.String(),
	}

	res := &pb.GetMenuByIDResponse{}
	useCase.GetMenuById(context.Background(), req, res)

	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusNotFound), res.ErrorCode)
	assert.Contains(t, res.ErrorMessage, mockMenuID.String())
	assert.Nil(t, res.Data)
}

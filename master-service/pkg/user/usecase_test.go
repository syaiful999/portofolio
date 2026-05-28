package user

import (
	"context"
	"testing"

	"moyo-master-service/config"
	pb "moyo-master-service/pkg/user/proto"
	"moyo-master-service/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockUserRepository is a mock implementation of IUserRepository for testing.
type MockUserRepository struct {
	MockGetUser            func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUser, int32, error)
	MockGetUserById        func(ctx context.Context, id uuid.UUID) (MasterUser, error)
	MockGetUserByEmail     func(ctx context.Context, email string) (MasterUser, string, error)
	MockCreateUser         func(ctx context.Context, value MasterUser) (MasterUser, string, error)
	MockUpdateUser         func(ctx context.Context, value MasterUser) (MasterUser, string, error)
	MockUpdateUserPassword func(ctx context.Context, value MasterUser, redirectId string) (MasterUser, error)

	MockDeleteUser                  func(ctx context.Context, id, modifiedBy uuid.UUID) error
	MockGetCountUniqueUser          func(ctx context.Context, value MasterUser) (int, error)
	MockGetUserGroupbyRole          func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUserGroupbyRole, error)
	MockActivateUser                func(ctx context.Context, value MasterUser) (MasterUser, error)
	MockGetDepartementByOutsourceId func(id string) ([]string, error)
	MockGetUserPasswordHash         func(ctx context.Context, userId uuid.UUID) (string, error)
}

func (m *MockUserRepository) GetUser(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUser, int32, error) {
	if m.MockGetUser != nil {
		return m.MockGetUser(ctx, skip, take, filter, sort)
	}
	return nil, 0, nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (MasterUser, string, error) {
	if m.MockGetUserByEmail != nil {
		return m.MockGetUserByEmail(ctx, email)
	}
	return MasterUser{}, "", nil
}
func (m *MockUserRepository) GetUserById(ctx context.Context, id uuid.UUID) (MasterUser, error) {
	if m.MockGetUserById != nil {
		return m.MockGetUserById(ctx, id)
	}
	return MasterUser{}, nil
}

func (m *MockUserRepository) CreateUser(ctx context.Context, value MasterUser) (MasterUser, string, error) {
	if m.MockCreateUser != nil {
		return m.MockCreateUser(ctx, value)
	}
	return MasterUser{}, "", nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, value MasterUser) (MasterUser, string, error) {
	if m.MockUpdateUser != nil {
		return m.MockUpdateUser(ctx, value)
	}
	return MasterUser{}, "", nil
}

func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, value MasterUser, redirectId string) (MasterUser, error) {
	if m.MockUpdateUserPassword != nil {
		return m.MockUpdateUserPassword(ctx, value, redirectId)
	}
	return MasterUser{}, nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id, modifiedBy uuid.UUID) error {
	if m.MockDeleteUser != nil {
		return m.MockDeleteUser(ctx, id, modifiedBy)
	}
	return nil
}
func (m *MockUserRepository) GetCountUniqueUser(ctx context.Context, value MasterUser) (int, error) {
	if m.MockGetCountUniqueUser != nil {
		return m.MockGetCountUniqueUser(ctx, value)
	}
	return 0, nil
}
func (m *MockUserRepository) GetUserGroupbyRole(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUserGroupbyRole, error) {
	if m.MockGetUserGroupbyRole != nil {
		return m.MockGetUserGroupbyRole(ctx, skip, take, filter, sort)
	}
	return []MasterUserGroupbyRole{}, nil
}
func (m *MockUserRepository) ActivateUser(ctx context.Context, value MasterUser) (MasterUser, error) {
	if m.MockGetUserGroupbyRole != nil {
		return m.MockActivateUser(ctx, value)
	}
	return value, nil
}

func (m *MockUserRepository) GetDepartementByOutsourceId(id string) ([]string, error) {
	if m.MockGetDepartementByOutsourceId != nil {
		return m.MockGetDepartementByOutsourceId(id)
	}
	return []string{}, nil
}

func (m *MockUserRepository) GetUserPasswordHash(ctx context.Context, userId uuid.UUID) (string, error) {
	if m.MockGetUserPasswordHash != nil {
		return m.MockGetUserPasswordHash(ctx, userId)
	}
	return "", nil
}

func TestGetUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var config config.Config
	useCase := NewUseCaseUser(mockRepo, config)

	// Mock implementation for GetUser
	mockRepo.MockGetUser = func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUser, int32, error) {
		return []MasterUser{{ID: uuid.New(), IsActive: true, CreatedBy: "user", UserName: "User1"}}, 1, nil
	}

	req := &pb.GetUserRequest{
		Skip:   0,
		Take:   1,
		Filter: "",
		Sort:   "",
	}

	token := &utils.TokenValue{}
	res := &pb.GetUserResponse{}
	err := useCase.GetUsers(context.Background(), req, res, token)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, int32(1), res.CountData)
}

/*

func TestUpdateUser(t *testing.T) {
	// Define a user to be updated
	idParse, _ := uuid.Parse("010680fb-4283-4c78-bfdc-36b5835b4e30")
	userToUpdate := MasterUser{
		ID:       idParse,
		Name:     "John Doe",
		Email:    "john@doe.com",
		IsActive: true,
	}

	// Instantiate the use case with the mock repository
	mockRepo := new(MockUserRepository)
	var config config.Config
	useCase := NewUseCaseUser(mockRepo, config)

	// Call the UpdateUser function
	err := usecase.UpdateUser(userToUpdate)

	// Check if the update was successful
	if err != nil {
		t.Errorf("Error updating user: %v", err)
	}

	// Check if the user was updated correctly
	if updatedUser.Id != userToUpdate.Id || updatedUser.Name != userToUpdate.Name || updatedUser.Email != userToUpdate.Email || updatedUser.IsActive != userToUpdate.IsActive {
		t.Errorf("UpdateUser did not update the user correctly. Expected: %v, Actual: %v", userToUpdate, updatedUser)
	}
}

// */

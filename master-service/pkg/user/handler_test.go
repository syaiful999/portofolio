package user

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"moyo-master-service/config"
	pb "moyo-master-service/pkg/user/proto"
	"moyo-master-service/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetUsers_NoToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)
	handler := NewUserHandler(useCase)

	req := &pb.GetUserRequest{Skip: 0, Take: 10}
	res := &pb.GetUserResponse{}

	err := handler.GetUsers(context.Background(), req, res)
	assert.Error(t, err)
}

func TestUserHandler_GetUserById_NoToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)
	handler := NewUserHandler(useCase)

	req := &pb.GetUserByIDRequest{Id: uuid.New().String()}
	res := &pb.GetUserByIDResponse{}

	err := handler.GetUserById(context.Background(), req, res)
	assert.Error(t, err)
}

func TestGetUsers_InvalidTake(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	req := &pb.GetUserRequest{Skip: 0, Take: 0}
	res := &pb.GetUserResponse{}
	token := &utils.TokenValue{ID: uuid.New()}

	err := useCase.GetUsers(context.Background(), req, res, token)
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusBadRequest), res.ErrorCode)
}

func TestGetUsers_RepositoryError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	userID := uuid.New()
	mockRepo.MockGetUserById = func(ctx context.Context, id uuid.UUID) (MasterUser, error) {
		return MasterUser{ID: userID, RoleCode: sql.NullString{String: "admin", Valid: true}}, nil
	}
	mockRepo.MockGetUser = func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUser, int32, error) {
		return nil, 0, errors.New("database error")
	}

	req := &pb.GetUserRequest{Skip: 0, Take: 10}
	res := &pb.GetUserResponse{}
	token := &utils.TokenValue{ID: userID}

	err := useCase.GetUsers(context.Background(), req, res, token)
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusInternalServerError), res.ErrorCode)
}

func TestGetUserById_InvalidUUID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	req := &pb.GetUserByIDRequest{Id: "invalid-uuid"}
	res := &pb.GetUserByIDResponse{}

	err := useCase.GetUserById(context.Background(), req, res)
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusBadRequest), res.ErrorCode)
}

func TestGetUserById_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	userID := uuid.New()
	mockRepo.MockGetUserById = func(ctx context.Context, id uuid.UUID) (MasterUser, error) {
		return MasterUser{
			ID:       userID,
			UserName: "testuser",
			Email:    "test@example.com",
			IsActive: true,
		}, nil
	}

	req := &pb.GetUserByIDRequest{Id: userID.String()}
	res := &pb.GetUserByIDResponse{}

	err := useCase.GetUserById(context.Background(), req, res)
	assert.NoError(t, err)
	assert.False(t, res.IsError)
	assert.NotNil(t, res.Data)
	assert.Equal(t, userID.String(), res.Data.Id)
	assert.Equal(t, "testuser", res.Data.UserName)
}

func TestGetUserById_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	userID := uuid.New()
	mockRepo.MockGetUserById = func(ctx context.Context, id uuid.UUID) (MasterUser, error) {
		return MasterUser{}, errors.New("sql: no rows in result set")
	}

	req := &pb.GetUserByIDRequest{Id: userID.String()}
	res := &pb.GetUserByIDResponse{}

	err := useCase.GetUserById(context.Background(), req, res)
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusNotFound), res.ErrorCode)
}

func TestCreateUser_MissingUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	req := &pb.MutationUserRequest{UserName: "", Email: "test@test.com"}
	res := &pb.MutationUserResponse{}
	token := &utils.TokenValue{ID: uuid.New()}

	err := useCase.CreateUser(context.Background(), req, res, token)
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusBadRequest), res.ErrorCode)
}

func TestCreateUser_DuplicateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	mockRepo.MockGetCountUniqueUser = func(ctx context.Context, value MasterUser) (int, error) {
		return 1, nil
	}

	req := &pb.MutationUserRequest{UserName: "existing", Email: "exist@test.com"}
	res := &pb.MutationUserResponse{}
	token := &utils.TokenValue{ID: uuid.New()}

	err := useCase.CreateUser(context.Background(), req, res, token)
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusBadRequest), res.ErrorCode)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	conf := config.Config{}
	conf.Hosts.Services.Web = "http://localhost:3000"
	useCase := NewUseCaseUser(mockRepo, conf)

	newID := uuid.New()
	mockRepo.MockGetCountUniqueUser = func(ctx context.Context, value MasterUser) (int, error) {
		return 0, nil
	}
	mockRepo.MockCreateUser = func(ctx context.Context, value MasterUser) (MasterUser, string, error) {
		return MasterUser{ID: newID}, "redirect-123", nil
	}

	req := &pb.MutationUserRequest{UserName: "newuser", Email: "new@test.com"}
	res := &pb.MutationUserResponse{}
	token := &utils.TokenValue{ID: uuid.New()}

	err := useCase.CreateUser(context.Background(), req, res, token)
	assert.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Equal(t, newID.String(), res.Data.Id)
}

func TestDeleteUser_InvalidUUID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	req := &pb.MutationUserRequest{Id: "invalid"}
	res := &pb.MutationUserResponse{}
	token := &utils.TokenValue{ID: uuid.New()}

	err := useCase.DeleteUser(context.Background(), req, res, token)
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusBadRequest), res.ErrorCode)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	userID := uuid.New()
	mockRepo.MockDeleteUser = func(ctx context.Context, id, modifiedBy uuid.UUID) error {
		return nil
	}

	req := &pb.MutationUserRequest{Id: userID.String()}
	res := &pb.MutationUserResponse{}
	token := &utils.TokenValue{ID: uuid.New()}

	err := useCase.DeleteUser(context.Background(), req, res, token)
	assert.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Equal(t, userID.String(), res.Data.Id)
}

func TestGetUserGroupbyRole_InvalidTake(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	req := &pb.GetUserRequest{Skip: 0, Take: 0}
	res := &pb.GetUserGroupByRoleResponse{}

	err := useCase.GetUserGroupbyRole(context.Background(), req, res)
	assert.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Equal(t, int32(http.StatusBadRequest), res.ErrorCode)
}

func TestGetUserGroupbyRole_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	var conf config.Config
	useCase := NewUseCaseUser(mockRepo, conf)

	mockRepo.MockGetUserGroupbyRole = func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUserGroupbyRole, error) {
		return []MasterUserGroupbyRole{
			{RoleCode: "admin", RoleDescription: "Administrator", Count: 5},
			{RoleCode: "user", RoleDescription: "Regular User", Count: 10},
		}, nil
	}

	req := &pb.GetUserRequest{Skip: 0, Take: 10}
	res := &pb.GetUserGroupByRoleResponse{}

	err := useCase.GetUserGroupbyRole(context.Background(), req, res)
	assert.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, "admin", res.Data[0].RoleCode)
	assert.Equal(t, int32(5), res.Data[0].Count)
}

package user

import (
	"context"
	"fmt"
	"moyo-master-service/config"
	pb "moyo-master-service/pkg/user/proto"
	"moyo-master-service/utils"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IUseCase interface {
	GetUsers(ctx context.Context, req *pb.GetUserRequest, res *pb.GetUserResponse, token *utils.TokenValue) error
	GetUserById(ctx context.Context, req *pb.GetUserByIDRequest, res *pb.GetUserByIDResponse) error
	UserForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest, res *pb.MutationUserResponse) error
	CreateUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse, token *utils.TokenValue) error
	UpdateUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse, token *utils.TokenValue) error
	UpdateUserPassword(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error
	DeleteUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse, token *utils.TokenValue) error
	GetUserGroupbyRole(ctx context.Context, req *pb.GetUserRequest, res *pb.GetUserGroupByRoleResponse) error
	ActivateUser(ctx context.Context, res *pb.MutationUserResponse, token *utils.TokenValue) error
}

type UseCase struct {
	repository IUserRepository
	config     config.Config
}

func NewUseCaseUser(repo IUserRepository, config config.Config) IUseCase {
	return &UseCase{
		repository: repo,
		config:     config,
	}
}

func errorResponse1(res *pb.GetUserResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse2(res *pb.GetUserByIDResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse3(res *pb.MutationUserResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse4(res *pb.GetUserGroupByRoleResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func (s *UseCase) GetUsers(ctx context.Context, req *pb.GetUserRequest, res *pb.GetUserResponse, token *utils.TokenValue) error {
	//Handle error param take

	users := make([]*pb.UserModel, 0)

	if req.Take <= 0 {
		return errorResponse1(
			res, http.StatusBadRequest, "GetUsers",
			fmt.Sprintf("invalid parameter take, Req.Take = %d", req.Take),
		)
	}

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "GetUsers",
			err.Error(),
		)
	}

	dataRequest, err := s.repository.GetUserById(ctx, token.ID)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "GetUsers:GetUserById",
			err.Error(),
		)
	}

	// role_code == 'admin_ocm' will only see users with the same role
	if strings.ToLower(dataRequest.RoleCode.String) == "admin_ocm" {
		if filter == "" {
			filter = fmt.Sprintf("WHERE role_id::TEXT = '%s'::TEXT", dataRequest.RoleId.String)
		} else {
			filter += fmt.Sprintf(" AND role_id::TEXT = '%s'::TEXT", dataRequest.RoleId.String)
		}
	}

	data, count, err := s.repository.GetUser(ctx, req.Skip, req.Take, filter, sort)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "GetUsers",
			err.Error(),
		)
	}

	for _, value := range data {
		var user pb.UserModel

		user.Id = value.ID.String()
		user.IsActive = value.IsActive
		user.IsStatus = value.IsStatus
		user.CreatedBy = value.CreatedBy
		user.CreatedDate = timestamppb.New(value.CreatedDate)
		user.ModifiedBy = value.ModifiedBy
		user.ModifiedDate = timestamppb.New(value.ModifiedDate)
		user.UserName = value.UserName
		user.Name = value.Name
		user.Email = value.Email
		user.RoleId = value.RoleId.String
		user.RoleDescription = value.RoleDescription.String
		user.RoleCode = value.RoleCode.String
		user.DepartmentId = value.DepartmentId.String
		user.DepartmentName = value.DepartmentName.String
		user.OutsourceId = value.OutsourceId.String
		user.OutsourceName = value.OutsourceName.String
		user.LocationId = value.LocationId.String
		user.LocationName = value.LocationName.String
		user.Picture = value.Picture.String
		users = append(users, &user)
	}
	res.Data = users
	res.CountData = count
	return nil
}

func (s *UseCase) GetUserById(ctx context.Context, req *pb.GetUserByIDRequest, res *pb.GetUserByIDResponse) error {
	var err error
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return errorResponse2(res, http.StatusBadRequest, "GetUserById:uuid", err.Error())
	}
	user, err := s.repository.GetUserById(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("user id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "GetUserById"+errStr, "")
		}
		return errorResponse2(res, http.StatusInternalServerError, "GetUserById", err.Error())
	}
	res.Data = &pb.UserModel{
		Id:              user.ID.String(),
		IsActive:        user.IsActive,
		IsStatus:        user.IsStatus,
		CreatedBy:       user.CreatedBy,
		CreatedDate:     timestamppb.New(user.CreatedDate),
		ModifiedBy:      user.ModifiedBy,
		ModifiedDate:    timestamppb.New(user.ModifiedDate),
		UserName:        user.UserName,
		Name:            user.Name,
		Email:           user.Email,
		RoleId:          user.RoleId.String,
		RoleDescription: user.RoleDescription.String,
		RoleCode:        user.RoleCode.String,
		DepartmentId:    user.DepartmentId.String,
		DepartmentName:  user.DepartmentName.String,
		OutsourceId:     user.OutsourceId.String,
		OutsourceName:   user.OutsourceName.String,
		LocationId:      user.LocationId.String,
		LocationName:    user.LocationName.String,
		Picture:         user.Picture.String,
	}

	return nil
}

func (s *UseCase) UserForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest, res *pb.MutationUserResponse) error {

	user, idRedirect, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("user email %s is not found", req.Email)
			return errorResponse3(res, http.StatusNotFound, "UserForgotPassword"+errStr, "")
		}
		return errorResponse3(res, http.StatusInternalServerError, "failed to UserForgotPassword", err.Error())
	}

	if user.IsStatus {
		linkActivation := fmt.Sprintf(`%s/reset-password?id=%s&userId=%s`, s.config.Hosts.Services.Web, idRedirect, user.ID.String())
		utils.PushLogf("", "UserForgotPassword", fmt.Sprintf("Reset password link: %s", linkActivation))
	} else {
		linkActivation := fmt.Sprintf(`%s/generate-password?id=%s&userId=%s`, s.config.Hosts.Services.Web, idRedirect, user.ID.String())
		utils.PushLogf("", "UserForgotPassword", fmt.Sprintf("Activate user link: %s", linkActivation))
	}

	return nil
}

func (s *UseCase) CreateUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse, token *utils.TokenValue) error {
	var err error
	var result MasterUser

	ID := uuid.New()

	//Handle mandatori field
	if utils.HandleMandatoryField(req.UserName) {
		return errorResponse3(res, http.StatusBadRequest, "CreateUser:user name is required", "")
	}

	data := MasterUser{
		ID:           ID,
		UserName:     req.UserName,
		Name:         req.Name,
		Email:        req.Email,
		RoleId:       utils.ConvStringToNullString(req.RoleId),
		DepartmentId: utils.ConvStringToNullString(req.DepartmentId),
		OutsourceId:  utils.ConvStringToNullString(req.OutsourceId),
		LocationId:   utils.ConvStringToNullString(req.LocationId),
		Picture:      utils.ConvStringToNullString(req.Picture),
		IsActive:     req.IsActive,
		IsStatus:     req.IsStatus,
		CreatedBy:    token.ID.String(),
		ModifiedBy:   token.ID.String(),
	}

	//check if exist
	countData, err := s.repository.GetCountUniqueUser(ctx, data)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "CreateUser:GetCountUniqueUser", err.Error())
	}
	if countData > 0 {
		errStr := fmt.Sprintf("email %s or username %s is already exist", req.Email, req.UserName)
		return errorResponse3(res, http.StatusBadRequest, "CreateUser"+errStr, "")
	}

	var idRedirect string
	if result, idRedirect, err = s.repository.CreateUser(ctx, data); err != nil {
		return errorResponse3(res, http.StatusBadRequest, "failed to CreateUser", err.Error())
	}

	linkActivation := fmt.Sprintf(`%s/generate-password?id=%s&userId=%s`, s.config.Hosts.Services.Web, idRedirect, result.ID.String())
	utils.PushLogf("", "CreateUser", fmt.Sprintf("Activate user link: %s", linkActivation))

	res.Data = &pb.GetUserByIDRequest{Id: result.ID.String()}
	return nil
}

func (s *UseCase) UpdateUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse, token *utils.TokenValue) error {
	var err error
	var result MasterUser
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "UpdateUser:uuid", err.Error())
	}

	if utils.HandleMandatoryField(req.UserName) {
		return errorResponse3(res, http.StatusBadRequest, "UpdateUser:user name is required", "")
	}

	data := MasterUser{
		ID:           ID,
		UserName:     req.UserName,
		Email:        req.Email,
		Name:         req.Name,
		RoleId:       utils.ConvStringToNullString(req.RoleId),
		DepartmentId: utils.ConvStringToNullString(req.DepartmentId),
		OutsourceId:  utils.ConvStringToNullString(req.OutsourceId),
		LocationId:   utils.ConvStringToNullString(req.LocationId),
		Picture:      utils.ConvStringToNullString(req.Picture),
		IsStatus:     req.IsStatus,
		IsActive:     req.IsActive,
		ModifiedBy:   token.ID.String(),
	}
	// check if exist
	countData, err := s.repository.GetCountUniqueUser(ctx, data)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "UpdateUser:GetCountUniqueUser", err.Error())
	}
	if countData > 0 {
		errStr := fmt.Sprintf("email %s or username %s is already exist", req.Email, req.UserName)
		return errorResponse3(res, http.StatusBadRequest, "UpdateUser"+errStr, "")
	}

	var idRedirect string
	if result, idRedirect, err = s.repository.UpdateUser(ctx, data); err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "failed to UpdateUser", err.Error())
	}
	if !result.IsStatus {
		linkActivation := fmt.Sprintf(`%s/generate-password?id=%s&userId=%s`, s.config.Hosts.Services.Web, idRedirect, result.ID.String())
		utils.PushLogf("", "UpdateUser", fmt.Sprintf("Activate user link: %s", linkActivation))
	}
	res.Data = &pb.GetUserByIDRequest{Id: result.ID.String()}
	return nil
}

func (s *UseCase) UpdateUserPassword(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error {
	var err error
	var result MasterUser

	if utils.HandleMandatoryField(req.Id) {
		return errorResponse3(res, http.StatusBadRequest, "UpdateUserPassword:id is required", "")
	}

	partsId := strings.Split(req.Id, "|")
	var userId uuid.UUID
	var redirectId string
	if len(partsId) > 1 {
		redirectId = partsId[0]
		userId, err = uuid.Parse(partsId[1])
		if err != nil {

			return errorResponse3(res, http.StatusBadRequest, "failed to UpdateUserPassword parseID", err.Error())
		}

	} else {
		return errorResponse3(res, http.StatusBadRequest, "failed to UpdateUserPassword:id is required", "")
	}

	// ✅ Step 1: Ambil hash password lama dari repository
	oldHashedPassword, err := s.repository.GetUserPasswordHash(ctx, userId)
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "UpdateUserPassword:failed to get user", "")
	}

	// 🔍 Debug log
	//utils.PushLogf("", "[DEBUG] UpdateUserPassword", "OldHash="+oldHashedPassword)

	// ✅ Step 2: Cek apakah password baru sama dengan lama
	if utils.CheckPasswordHash(req.Password, oldHashedPassword) {
		return errorResponse3(res, http.StatusBadRequest, "UpdateUserPassword:new password must be different from old password", "")
	}

	// ✅ Step 3: Hash password baru
	encriptedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "UpdateUserPassword:error hashing password", "")
	}

	// ✅ Step 4: Update password ke DB
	data := MasterUser{
		ID:         userId,
		Password:   encriptedPassword,
		ModifiedBy: userId.String(),
	}

	if result, err = s.repository.UpdateUserPassword(ctx, data, redirectId); err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "UpdateUserPassword", err.Error())
	}
	res.Data = &pb.GetUserByIDRequest{Id: result.ID.String()}
	return nil
}
func (s *UseCase) DeleteUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse, token *utils.TokenValue) error {
	var err error

	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "failed to DeleteUser:uuid", err.Error())
	}
	err = s.repository.DeleteUser(ctx, ID, token.ID)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "failed to DeleteUser", err.Error())
	}

	res.Data = &pb.GetUserByIDRequest{Id: ID.String()}
	return nil
}

func (s *UseCase) GetUserGroupbyRole(ctx context.Context, req *pb.GetUserRequest, res *pb.GetUserGroupByRoleResponse) error {

	usersGroupbyRole := make([]*pb.UserRoleModel, 0)
	if req.Take <= 0 {
		return errorResponse4(
			res, http.StatusBadRequest, "GetUserGroupbyRole",
			fmt.Sprintf("invalid parameter take, Req.Take = %d", req.Take),
		)
	}

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)

	if err != nil {
		return errorResponse4(
			res, http.StatusInternalServerError, "failed to GetUserGroupbyRole:KendoSortFilter",
			err.Error(),
		)
	}

	data, err := s.repository.GetUserGroupbyRole(ctx, req.Skip, req.Take, filter, sort)
	if err != nil {
		return errorResponse4(
			res, http.StatusInternalServerError, "failed to GetUserGroupbyRole",
			err.Error(),
		)
	}

	for _, value := range data {
		var userGrouped pb.UserRoleModel

		userGrouped.Count = value.Count
		userGrouped.RoleCode = value.RoleCode
		userGrouped.RoleDescription = value.RoleDescription

		usersGroupbyRole = append(usersGroupbyRole, &userGrouped)
	}
	res.Data = usersGroupbyRole
	return nil
}

func (s *UseCase) ActivateUser(ctx context.Context, res *pb.MutationUserResponse, token *utils.TokenValue) error {
	var result MasterUser
	var err error

	data := MasterUser{
		ID: token.ID,
	}

	if result, err = s.repository.ActivateUser(ctx, data); err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "failed to ActivateUser", err.Error())
	}
	res.Data = &pb.GetUserByIDRequest{Id: result.ID.String()}
	return nil
}

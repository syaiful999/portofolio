package user

import (
	"context"
	pb "moyo-master-service/pkg/user/proto"
	"moyo-master-service/utils"
)

type UserHandler struct {
	srv IUseCase
}

type IUserHandler interface {
	GetUsers(ctx context.Context, req *pb.GetUserRequest, res *pb.GetUserResponse) error
	GetUserById(ctx context.Context, req *pb.GetUserByIDRequest, res *pb.GetUserByIDResponse) error
	CreateUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error
	UpdateUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error
	UpdateUserPassword(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error
	DeleteUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error
	GetUserGroupbyRole(ctx context.Context, req *pb.GetUserRequest, res *pb.GetUserGroupByRoleResponse) error
	ActivateUser(ctx context.Context, req *pb.Empty, res *pb.MutationUserResponse) error
	UserForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest, res *pb.MutationUserResponse) error
}

func NewUserHandler(srv IUseCase) IUserHandler {
	return &UserHandler{
		srv: srv,
	}
}

func (e *UserHandler) GetUsers(ctx context.Context, req *pb.GetUserRequest, res *pb.GetUserResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetUsers(ctx, req, res, token)
}

func (e *UserHandler) GetUserById(ctx context.Context, req *pb.GetUserByIDRequest, res *pb.GetUserByIDResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetUserById(ctx, req, res)
}

func (e *UserHandler) CreateUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.CreateUser(ctx, req, res, token)
}

func (e *UserHandler) UpdateUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateUser(ctx, req, res, token)
}
func (e *UserHandler) UpdateUserPassword(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error {
	return e.srv.UpdateUserPassword(ctx, req, res)
}

func (e *UserHandler) DeleteUser(ctx context.Context, req *pb.MutationUserRequest, res *pb.MutationUserResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.DeleteUser(ctx, req, res, token)
}

func (e *UserHandler) GetUserGroupbyRole(ctx context.Context, req *pb.GetUserRequest, res *pb.GetUserGroupByRoleResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetUserGroupbyRole(ctx, req, res)
}

func (e *UserHandler) ActivateUser(ctx context.Context, req *pb.Empty, res *pb.MutationUserResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.ActivateUser(ctx, res, token)
}
func (e *UserHandler) UserForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest, res *pb.MutationUserResponse) error {
	return e.srv.UserForgotPassword(ctx, req, res)
}

package enum

import (
	"context"
	pb "moyo-master-service/pkg/enum/proto"
	"moyo-master-service/utils"
)

type EnumHandler struct {
	srv IUseCase
}

type IEnumHandler interface {
	GetEnums(ctx context.Context, req *pb.GetEnumRequest, res *pb.GetEnumResponse) error
	GetEnumById(ctx context.Context, req *pb.GetEnumByIDRequest, res *pb.GetEnumByIDResponse) error
	UpdateEnum(ctx context.Context, req *pb.UpdateEnumRequest, res *pb.UpdateEnumResponse) error
}

func NewEnumHandler(srv IUseCase) IEnumHandler {
	return &EnumHandler{
		srv: srv,
	}
}

func (e *EnumHandler) GetEnums(ctx context.Context, req *pb.GetEnumRequest, res *pb.GetEnumResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEnums(ctx, req, res)
}

func (e *EnumHandler) GetEnumById(ctx context.Context, req *pb.GetEnumByIDRequest, res *pb.GetEnumByIDResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEnumById(ctx, req, res)
}

func (e *EnumHandler) UpdateEnum(ctx context.Context, req *pb.UpdateEnumRequest, res *pb.UpdateEnumResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEnum(ctx, req, res, token)
}

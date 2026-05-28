package menu

import (
	"context"
	pb "moyo-master-service/pkg/menu/proto"
	"moyo-master-service/utils"
)

type MenuHandler struct {
	srv IUseCase
}

type IMenuHandler interface {
	GetMenu(ctx context.Context, req *pb.GetMenuRequest, res *pb.GetMenuResponse) error
	GetMenuByID(ctx context.Context, req *pb.GetMenuByIDRequest, res *pb.GetMenuByIDResponse) error
}

func NewMenuHandler(srv IUseCase) IMenuHandler {
	return &MenuHandler{
		srv: srv,
	}
}

func (e *MenuHandler) GetMenu(ctx context.Context, req *pb.GetMenuRequest, res *pb.GetMenuResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetMenus(ctx, req, res)
}

func (e *MenuHandler) GetMenuByID(ctx context.Context, req *pb.GetMenuByIDRequest, res *pb.GetMenuByIDResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetMenuById(ctx, req, res)
}

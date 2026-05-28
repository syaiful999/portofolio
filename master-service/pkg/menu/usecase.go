package menu

import (
	"context"
	"fmt"
	pb "moyo-master-service/pkg/menu/proto"
	"moyo-master-service/utils"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IUseCase interface {
	GetMenus(ctx context.Context, req *pb.GetMenuRequest, res *pb.GetMenuResponse) error
	GetMenuById(ctx context.Context, req *pb.GetMenuByIDRequest, res *pb.GetMenuByIDResponse) error
}

type UseCase struct {
	repository IMenuRepository
}

func NewUseCaseMenu(repo IMenuRepository) IUseCase {
	return &UseCase{
		repository: repo,
	}
}

func errorResponse1(res *pb.GetMenuResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse2(res *pb.GetMenuByIDResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func (s *UseCase) GetMenus(ctx context.Context, req *pb.GetMenuRequest, res *pb.GetMenuResponse) error {
	//Handle error param take

	if req.Take <= 0 {
		return errorResponse1(
			res, http.StatusBadRequest,
			fmt.Sprintf("GetMenus:invalid parameter take, Req.Take = %d", req.Take),
			"",
		)
	}

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "GetMenus:KendoSortFilter",
			err.Error(),
		)
	}

	data, count, err := s.repository.GetMenu(ctx, req.Skip, req.Take, filter, sort)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "failed to GetMenus:",
			err.Error(),
		)
	}

	var menus []*pb.MenuModel
	flatMenus := make([]*pb.MenuModel, 0)

	for _, value := range data {
		var rfid pb.MenuModel

		rfid.Id = value.ID.String()
		rfid.MenuCode = value.MenuCode
		rfid.MenuDescription = value.MenuDescription
		rfid.MenuPath = value.MenuPath
		rfid.IsActive = value.IsActive
		rfid.CreatedBy = value.CreatedBy
		rfid.CreatedDate = timestamppb.New(value.CreatedDate)
		rfid.ModifiedBy = value.ModifiedBy
		rfid.ModifiedDate = timestamppb.New(value.ModifiedDate)
		rfid.IdParent = value.IdParent.UUID.String()
		rfid.Level = value.Level

		flatMenus = append(flatMenus, &rfid)
	}
	menus = BuildMenuTree(flatMenus)

	res.Data = menus
	res.CountData = count
	return nil
}

func (s *UseCase) GetMenuById(ctx context.Context, req *pb.GetMenuByIDRequest, res *pb.GetMenuByIDResponse) error {
	var err error
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return errorResponse2(res, http.StatusBadRequest, "GetMenuById:uuid", err.Error())
	}
	menu, err := s.repository.GetMenuById(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("menu id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "GetMenuById"+errStr, "")
		}
		return errorResponse2(res, http.StatusInternalServerError, "failed to GetMenuById", err.Error())
	}

	var menus []*pb.MenuModel
	flatMenus := make([]*pb.MenuModel, 0)

	for _, value := range menu {
		var rfid pb.MenuModel

		rfid.Id = value.ID.String()
		rfid.MenuCode = value.MenuCode
		rfid.MenuDescription = value.MenuDescription
		rfid.MenuPath = value.MenuPath
		rfid.IsActive = value.IsActive
		rfid.CreatedBy = value.CreatedBy
		rfid.CreatedDate = timestamppb.New(value.CreatedDate)
		rfid.ModifiedBy = value.ModifiedBy
		rfid.ModifiedDate = timestamppb.New(value.ModifiedDate)
		rfid.IdParent = value.IdParent.UUID.String()
		rfid.Level = value.Level

		flatMenus = append(flatMenus, &rfid)
	}
	menus = BuildMenuTree(flatMenus)

	res.Data = menus[0]
	return nil
}

func BuildMenuTree(flat []*pb.MenuModel) []*pb.MenuModel {
	menuMap := make(map[string]*pb.MenuModel)
	var roots []*pb.MenuModel
	for _, m := range flat {
		m.Children = []*pb.MenuModel{}
		menuMap[m.Id] = m
	}
	for _, m := range flat {

		if m.IdParent == "" || m.IdParent == "0" || menuMap[m.IdParent] == nil {
			roots = append(roots, m)
			continue
		}
		parent := menuMap[m.IdParent]
		parent.Children = append(parent.Children, m)
	}

	return roots
}

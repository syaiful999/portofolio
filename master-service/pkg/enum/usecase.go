package enum

import (
	"context"
	"fmt"
	pb "moyo-master-service/pkg/enum/proto"
	"moyo-master-service/utils"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IUseCase interface {
	GetEnums(ctx context.Context, req *pb.GetEnumRequest, res *pb.GetEnumResponse) error
	GetEnumById(ctx context.Context, req *pb.GetEnumByIDRequest, res *pb.GetEnumByIDResponse) error
	UpdateEnum(ctx context.Context, req *pb.UpdateEnumRequest, res *pb.UpdateEnumResponse, token *utils.TokenValue) error
}

type UseCase struct {
	repository IEnumRepository
}

func NewUseCaseEnum(repo IEnumRepository) IUseCase {
	return &UseCase{
		repository: repo,
	}
}

func errorResponse1(res *pb.GetEnumResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse2(res *pb.GetEnumByIDResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse3(res *pb.UpdateEnumResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func (s *UseCase) GetEnums(ctx context.Context, req *pb.GetEnumRequest, res *pb.GetEnumResponse) error {
	//Handle error param take

	enums := make([]*pb.EnumModel, 0)

	if req.Take <= 0 {
		return errorResponse1(
			res, http.StatusBadRequest, "GetEnums:",
			fmt.Sprintf("invalid parameter take, Req.Take = %d", req.Take),
		)
	}

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "GetEnums:kendo",
			err.Error(),
		)
	}

	data, count, err := s.repository.GetEnum(ctx, req.Skip, req.Take, filter, sort)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "GetEnums:",
			err.Error(),
		)
	}

	for _, value := range data {
		var enum pb.EnumModel

		enum.Id = value.ID.String()
		enum.EnumValue = value.EnumValue
		enum.EnumType = value.EnumType
		enum.EnumCode = value.EnumCode.String
		enum.IsActive = value.IsActive
		enum.CreatedBy = value.CreatedBy
		enum.CreatedDate = timestamppb.New(value.CreatedDate)
		enum.ModifiedBy = value.ModifiedBy
		enum.ModifiedDate = timestamppb.New(value.ModifiedDate)

		enums = append(enums, &enum)
	}
	res.Data = enums
	res.CountData = count
	return nil
}

func (s *UseCase) GetEnumById(ctx context.Context, req *pb.GetEnumByIDRequest, res *pb.GetEnumByIDResponse) error {
	var err error
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return errorResponse2(res, http.StatusBadRequest, "failed to GetEnumById:uuid", err.Error())
	}
	enum, err := s.repository.GetEnumById(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("enum id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "failed to GetEnumById"+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "failed to GetEnumById", err.Error())
	}
	res.Data = &pb.EnumModel{
		Id:           enum.ID.String(),
		EnumValue:    enum.EnumValue,
		EnumType:     enum.EnumType,
		IsActive:     enum.IsActive,
		CreatedBy:    enum.CreatedBy,
		CreatedDate:  timestamppb.New(enum.CreatedDate),
		ModifiedBy:   enum.ModifiedBy,
		ModifiedDate: timestamppb.New(enum.ModifiedDate),
	}

	return nil
}

func (s *UseCase) UpdateEnum(ctx context.Context, req *pb.UpdateEnumRequest, res *pb.UpdateEnumResponse, token *utils.TokenValue) error {
	var err error
	var result MasterEnum
	if utils.HandleMandatoryField(req.EnumValue) {
		res.Data = &pb.GetEnumByIDRequest{Id: result.ID.String()}
		return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEnum: enum value is required", "")
	}

	data := MasterEnum{
		EnumType:   req.EnumType,
		EnumValue:  req.EnumValue,
		ModifiedBy: token.ID.String(),
	}

	//update data
	if result, err = s.repository.UpdateEnum(ctx, data); err != nil {
		res.Data = &pb.GetEnumByIDRequest{Id: result.ID.String()}
		return errorResponse3(res, http.StatusInternalServerError, "failed to UpdateEnum:", err.Error())
	}

	res.Data = &pb.GetEnumByIDRequest{Id: result.ID.String()}
	return nil
}

package enum

import (
	"context"
	"testing"

	pb "moyo-master-service/pkg/enum/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockUseCase struct {
	MockGetEnums    func(ctx context.Context, req *pb.GetEnumRequest, res *pb.GetEnumResponse) error
	MockGetEnumById func(ctx context.Context, req *pb.GetEnumByIDRequest, res *pb.GetEnumByIDResponse) error
	MockUpdateEnum  func(ctx context.Context, req *pb.UpdateEnumRequest, res *pb.UpdateEnumResponse, token interface{}) error
}

func TestEnumHandler_GetEnums_NoToken(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)
	handler := NewEnumHandler(useCase)

	req := &pb.GetEnumRequest{Skip: 0, Take: 10}
	res := &pb.GetEnumResponse{}

	err := handler.GetEnums(context.Background(), req, res)
	assert.Error(t, err)
}

func TestEnumHandler_GetEnumById_NoToken(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)
	handler := NewEnumHandler(useCase)

	req := &pb.GetEnumByIDRequest{Id: uuid.New().String()}
	res := &pb.GetEnumByIDResponse{}

	err := handler.GetEnumById(context.Background(), req, res)
	assert.Error(t, err)
}

func TestEnumHandler_UpdateEnum_NoToken(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	useCase := NewUseCaseEnum(mockRepo)
	handler := NewEnumHandler(useCase)

	req := &pb.UpdateEnumRequest{Id: uuid.New().String(), EnumValue: "test"}
	res := &pb.UpdateEnumResponse{}

	err := handler.UpdateEnum(context.Background(), req, res)
	assert.Error(t, err)
}

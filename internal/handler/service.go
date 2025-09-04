package handler

import (
	"context"
	"fmt"

	"github.com/luzmareto/go-grpc-ecommerce-be/internal/utils"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/service"
)

type serivceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (sh *serivceHandler) HelloWorld(ctx context.Context, request *service.HelloWordlRequest) (*service.HelloWorldResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &service.HelloWorldResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("hello %s", request.Name),
		Base:    utils.SuccessResponse("Success"),
	}, nil

}

func NewServiceHandler() *serivceHandler {
	return &serivceHandler{}
}

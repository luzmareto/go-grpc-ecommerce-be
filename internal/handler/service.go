package handler

import (
	"context"
	"fmt"

	"github.com/luzmareto/go-grpc-ecommerce-be/pb/service"
)

type serivceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (sh *serivceHandler) HelloWorld(ctx context.Context, request *service.HelloWordlRequest) (*service.HelloWorldResponse, error) {
	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s", request.Name),
	}, nil
}

func NewServiceHandler() *serivceHandler {
	return &serivceHandler{}
}

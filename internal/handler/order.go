package handler

import (
	"context"

	"github.com/luzmareto/go-grpc-ecommerce-be/internal/service"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/utils"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/order"
)

type orderHandler struct {
	order.UnimplementedOrderServiceServer

	orderService service.IOrderService
}

func (oh *orderHandler) CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &order.CreateOrderResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.CreateOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewOrderHandler(orderService service.IOrderService) *orderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

package service

import (
	"context"
	"fmt"
	"time"

	jwtentity "github.com/luzmareto/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/repository"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/utils"

	"github.com/google/uuid"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/entity"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/order"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
}

type orderService struct {
	orderRepository  repository.IOrderRepository
	productRepostory repository.IProductRepository
}

func (os *orderService) CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	numbering, err := os.orderRepository.GetNumbering(ctx, "order")
	if err != nil {
		return nil, err
	}

	var productIds = make([]string, len(request.Products))
	for i := range request.Products {
		productIds[i] = request.Products[i].Id
	}

	products, err := os.productRepostory.GetProductsByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]*entity.Product)
	for i := range products {
		productMap[products[i].Id] = products[i]
	}

	var total float64 = 0
	for _, p := range request.Products {
		if productMap[p.Id] == nil {
			return &order.CreateOrderResponse{
				Base: utils.NotFoundResponse(fmt.Sprintf("Product %s not found", p.Id)),
			}, nil
		}
		total += productMap[p.Id].Price * float64(p.Quantity)
	}
	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)
	orderEntity := entity.Order{
		Id:              uuid.NewString(),
		Number:          fmt.Sprintf("ORD-%d%08d", now.Year(), numbering.Number),
		UserId:          claims.Subject,
		OrderStatusCode: entity.OrderStatusCodeUnpaid,
		UserFullName:    request.FullName,
		Address:         request.Address,
		PhoneNumber:     request.PhoneNumber,
		Notes:           &request.Notes,
		Total:           total,
		ExpiredAt:       &expiredAt,
		CreatedAt:       now,
		CreatedBy:       claims.FullName,
	}

	err = os.orderRepository.CreateOrder(ctx, &orderEntity)
	if err != nil {
		return nil, err
	}

	for _, p := range request.Products {
		var orderItem = entity.OrderItem{
			Id:                   uuid.NewString(),
			ProductId:            p.Id,
			ProductName:          productMap[p.Id].Name,
			ProductImageFileName: productMap[p.Id].ImageFileName,
			ProductPrice:         productMap[p.Id].Price,
			Quantity:             p.Quantity,
			OrderId:              orderEntity.Id,
			CreatedAt:            now,
			CreatedBy:            claims.FullName,
		}

		err = os.orderRepository.CreateOrderItem(ctx, &orderItem)
		if err != nil {
			return nil, err
		}
	}

	numbering.Number++
	err = os.orderRepository.UpdateNumbering(ctx, numbering)
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		Base: utils.SuccessResponse("Create order success"),
		Id:   orderEntity.Id,
	}, nil
}

// func (os *orderService) CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
// 	claims, err := jwtentity.GetClaimsFromContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	numbering, err := os.orderRepository.GetNumbering(ctx, "order")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var productIds = make([]string, len(request.Products))
// 	for i := range request.Products {
// 		productIds[i] = request.Products[i].Id
// 	}

// 	products, err := os.productRepostory.GetProductsByIds(ctx, productIds)
// 	if err != nil {
// 		return nil, err
// 	}

// 	productMap := make(map[string]*entity.Product)
// 	for i := range products {
// 		productMap[products[i].Id] = products[i]
// 	}

// 	// ✅ 1️⃣ Validasi SEMUA produk sebelum insert apapun ke DB
// 	for _, p := range request.Products {
// 		if productMap[p.Id] == nil {
// 			return &order.CreateOrderResponse{
// 				Base: utils.NotFoundResponse(fmt.Sprintf("Product with id %s not found", p.Id)),
// 			}, nil
// 		}
// 	}

// 	// ✅ 2️⃣ Kalau semua valid, baru lanjut
// 	var total float64
// 	for _, p := range request.Products {
// 		total += productMap[p.Id].Price * float64(p.Quantity)
// 	}

// 	now := time.Now()
// 	expiredAt := now.Add(24 * time.Hour)
// 	orderEntity := entity.Order{
// 		Id:              uuid.NewString(),
// 		Number:          fmt.Sprintf("ORD-%d%08d", now.Year(), numbering.Number),
// 		UserId:          claims.Subject,
// 		OrderStatusCode: entity.OrderStatusCodeUnpaid,
// 		UserFullName:    request.FullName,
// 		Address:         request.Address,
// 		PhoneNumber:     request.PhoneNumber,
// 		Notes:           &request.Notes,
// 		Total:           total,
// 		ExpiredAt:       &expiredAt,
// 		CreatedAt:       now,
// 		CreatedBy:       claims.FullName,
// 	}

// 	err = os.orderRepository.CreateOrder(ctx, &orderEntity)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, p := range request.Products {
// 		orderItem := entity.OrderItem{
// 			Id:                   uuid.NewString(),
// 			ProductId:            p.Id,
// 			ProductName:          productMap[p.Id].Name,
// 			ProductImageFileName: productMap[p.Id].ImageFileName,
// 			ProductPrice:         productMap[p.Id].Price,
// 			Quantity:             p.Quantity,
// 			OrderId:              orderEntity.Id,
// 			CreatedAt:            now,
// 			CreatedBy:            claims.FullName,
// 		}
// 		err = os.orderRepository.CreateOrderItem(ctx, &orderItem)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	numbering.Number++
// 	err = os.orderRepository.UpdateNumbering(ctx, numbering)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &order.CreateOrderResponse{
// 		Base: utils.SuccessResponse("Create order success"),
// 		Id:   orderEntity.Id,
// 	}, nil
// }

func NewOrderService(orderRepository repository.IOrderRepository, productRepository repository.IProductRepository) IOrderService {
	return &orderService{
		orderRepository:  orderRepository,
		productRepostory: productRepository,
	}
}

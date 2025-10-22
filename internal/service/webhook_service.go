package service

import (
	"context"
	"errors"
	"time"

	"github.com/luzmareto/go-grpc-ecommerce-be/internal/dto"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/entity"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/repository"
)

type IWebhookService interface {
	ReceiveInvoice(ctx context.Context, request *dto.XenditInvoiceRequest) error
}

type webhookService struct {
	orderRepository repository.IOrderRepository
}

func (ws *webhookService) ReceiveInvoice(ctx context.Context, request *dto.XenditInvoiceRequest) error {
	orderEntity, err := ws.orderRepository.GetOrderById(ctx, request.ExternalID)
	if err != nil {
		return err
	} //logic jika order tidak ditemukan
	if orderEntity == nil {
		return errors.New("order not found")
	}

	now := time.Now()
	updatedBy := "System"
	orderEntity.OrderStatusCode = entity.OrderStatusCodePaid
	orderEntity.UpdatedAt = &now
	orderEntity.UpdatedBy = &updatedBy
	orderEntity.XenditPaidAt = &now
	orderEntity.XenditPaymentChannel = &request.PaymentChannel
	orderEntity.XenditPaymentMethod = &request.PaymentMethod

	err = ws.orderRepository.UpdateOrder(ctx, orderEntity)
	if err != nil {
		return err
	}

	return nil
}

func NewWebhookService(orderRepository repository.IOrderRepository) IWebhookService {
	return &webhookService{
		orderRepository: orderRepository,
	}
}

package handler

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/dto"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/service"
)

type webhookHandler struct {
	webhookService service.IWebhookService
}

func (wh *webhookHandler) ReceiveInvoice(c *fiber.Ctx) error {
	var request dto.XenditInvoiceRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusBadRequest)
	}

	err = wh.webhookService.ReceiveInvoice(c.UserContext(), &request)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

func NewWebhookHandler(webhookService service.IWebhookService) *webhookHandler {
	return &webhookHandler{
		webhookService: webhookService,
	}
}

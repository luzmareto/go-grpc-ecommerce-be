package handler

import (
	"context"

	"github.com/luzmareto/go-grpc-ecommerce-be/internal/service"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/utils"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/newsletter"
)

type newsletterHandler struct {
	newsletter.UnimplementedNewsletterServiceServer

	newsletterService service.InewsLetterService
}

func (nh *newsletterHandler) SubcribeNewsletter(ctx context.Context, request *newsletter.SubcribeNewsletterRequest) (*newsletter.SubcribeNewsletterResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &newsletter.SubcribeNewsletterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := nh.newsletterService.SubcribeNewsletter(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func NewNewsletterHandler(newsletterSerive service.InewsLetterService) *newsletterHandler {
	return &newsletterHandler{
		newsletterService: newsletterSerive,
	}
}

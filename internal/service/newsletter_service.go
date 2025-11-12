package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/entity"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/repository"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/utils"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/newsletter"
)

type InewsLetterService interface {
	SubcribeNewsletter(ctx context.Context, request *newsletter.SubcribeNewsletterRequest) (*newsletter.SubcribeNewsletterResponse, error)
}

type newsletterService struct {
	newsletterRepository repository.InewsLetterRepository
}

func (ns *newsletterService) SubcribeNewsletter(ctx context.Context, request *newsletter.SubcribeNewsletterRequest) (*newsletter.SubcribeNewsletterResponse, error) {

	newsletterEntity, err := ns.newsletterRepository.GetNewsLetterByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if newsletterEntity != nil {
		return &newsletter.SubcribeNewsletterResponse{
			Base: utils.SuccessResponse("Subcribe newsletter success"),
		}, nil
	}

	newNewsletterEntity := entity.Newsletter{ // insert db
		Id:        uuid.NewString(),
		Fullname:  request.FullName,
		Email:     request.Email,
		CreatedAt: time.Now(),
		CreatedBy: "Public",
	}

	err = ns.newsletterRepository.CreateNewNewsletter(ctx, &newNewsletterEntity)
	if err != nil {
		return nil, err
	}
	return &newsletter.SubcribeNewsletterResponse{
		Base: utils.SuccessResponse("Subcribe newsletter success"),
	}, nil
}

func NewNewsLetterService(newsletterRepository repository.InewsLetterRepository) InewsLetterService {
	return &newsletterService{
		newsletterRepository: newsletterRepository,
	}
}

package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/entity"
	jwtentity "github.com/luzmareto/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/repository"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/utils"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/product"
)

type IProductService interface {
	CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error)
	DetailProduct(ctx context.Context,request *product.DetailProductRequest) (*product.DetailProductResponse, error)
}

type productService struct {
	productRepository repository.IProductRepository
}

func (ps *productService) CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// hanya admin yang bisa create product
	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	imagePath := filepath.Join("storage", "product",request.ImageFileName)
	_, err = os.Stat(imagePath)
	if err != nil {
		if os.IsNotExist(err){
			return  &product.CreateProductResponse{
				Base: utils.BadRequestResponse("file not found"),
			}, nil
		}

		return nil, err
	}

	productEntity := entity.Product{
		Id:            uuid.NewString(),
		Name:          request.Name,
		Description:   request.Description,
		Price:         request.Price,
		ImageFileName: request.ImageFileName,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.FullName,
	}
	err = ps.productRepository.CreateNewProduct(ctx, &productEntity)
	if err != nil {
		return nil, err
	}

	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("Product is created"),
		Id:   productEntity.Id,
	}, nil
}

func (ps *productService) DetailProduct(ctx context.Context,request *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	productEntity, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if productEntity == nil {
		return &product.DetailProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}
	return &product.DetailProductResponse{
		Base: utils.SuccessResponse("Get product detail success"),
		Id: productEntity.Id,
		Name: productEntity.Name,
		Description: productEntity.Description,
		Price: productEntity.Price,
		ImageUrl: fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), productEntity.ImageFileName),
	}, nil
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}

package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/entity"
	jwtentity "github.com/luzmareto/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/repository"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/utils"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/cart"
)

type ICartService interface {
	AddProductToCart(ctx context.Context, request *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error)
	ListCart(ctx context.Context,request *cart.ListCartRequest) (*cart.ListCartResponse, error)
	DeleteCart(ctx context.Context,request *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error)
	UpdateCartQuantity(ctx context.Context, request *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error)
}

type cartService struct {
	productRepository repository.IProductRepository
	cartRepository repository.ICartRepository
}

func (cs *cartService) AddProductToCart(ctx context.Context, request *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	
	productEntity, err := cs.productRepository.GetProductById(ctx,request.ProductId)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &cart.AddProductToCartResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	cartEntity, err := cs.cartRepository.GetCartByProductAndUserId(ctx, request.ProductId, claims.Subject)
	if err != nil {
		return nil, err
	}

	if cartEntity != nil {
		now := time.Now()
		cartEntity.Quantity += 1
		cartEntity.UpdatedAt = &now
		cartEntity.UpdatedBy = &claims.Subject

		err = cs.cartRepository.UpdateCart(ctx, cartEntity)
		if err != nil {
			return nil, err
		}

		return &cart.AddProductToCartResponse{
		Base: utils.SuccessResponse(" Add product to cart success"),
		Id: cartEntity.Id,
	},nil 
	}

	newCartEntity := entity.UserCart{
		Id: uuid.NewString(),
		UserId: claims.Subject,
		ProductId: request.ProductId,
		Quantity: 1,
		CreatedAt: time.Now(),
		CreatedBy: claims.FullName,
	}
	err = cs.cartRepository.CreateNewCart(ctx, &newCartEntity)
	if err != nil {
		return nil, err
	}

	return &cart.AddProductToCartResponse{
		Base: utils.SuccessResponse(" Add product to cart success"),
		Id: newCartEntity.Id,
	},nil
}

func (cs *cartService) ListCart(ctx context.Context,request *cart.ListCartRequest) (*cart.ListCartResponse, error){
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	carts, err :=  cs.cartRepository.GetlistCart(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	var items []*cart.ListCartResponseItem = make([]*cart.ListCartResponseItem, 0)
	for _, cartEntity := range carts {
		item := cart.ListCartResponseItem{
			CartId: cartEntity.Id,
			ProductId: cartEntity.Product.Id,
			ProductName: cartEntity.Product.Name,
			ProductImageUrl: fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), cartEntity.Product.ImageFileName),
			ProductPrice: cartEntity.Product.Price,
			Quantity: int64(cartEntity.Quantity),
		}
	
	items = append(items, &item)
	}

	return &cart.ListCartResponse{
		Base: utils.SuccessResponse("Get list cart success"),
		Items: items,
	},nil
	
}

func (cs *cartService) DeleteCart(ctx context.Context,request *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error){
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cartEtnity, err := cs.cartRepository.GetCartById(ctx, request.CartId)
	if err != nil {
		return nil, err
	}
	if cartEtnity == nil {
		return  &cart.DeleteCartResponse{
			Base: utils.NotFoundResponse("Cart not found"),
		}, nil
	}
	
	if cartEtnity.UserId != claims.Subject{
		return &cart.DeleteCartResponse{
			Base: utils.BadRequestResponse("Cart user is not matched"),
		}, nil
	}

	err  = cs.cartRepository.DeleteCart(ctx, request.CartId)
	if err != nil {
		return nil, err
	}
	return  &cart.DeleteCartResponse{
		Base: utils.SuccessResponse("Delete cart success"),
	}, nil
}

func (cs *cartService) UpdateCartQuantity(ctx context.Context, request *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error){
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	
	cartEntity, err := cs.cartRepository.GetCartById(ctx, request.CartId)
	if err != nil {
		return nil, err
	}
	if cartEntity == nil {
		return  &cart.UpdateCartQuantityResponse{
			Base: utils.NotFoundResponse("Cart not found"),
		}, nil
	}

	if cartEntity.UserId != claims.Subject {
		return  &cart.UpdateCartQuantityResponse{
			Base: utils.BadRequestResponse("Cart user is not matched"),
		}, nil		
	}

	if request.NewQuantity == 0 {
		err =cs.cartRepository.DeleteCart(ctx, request.CartId)
		if err != nil {
			return nil, err
		}

		return &cart.UpdateCartQuantityResponse{
			Base: utils.SuccessResponse("Update cart quantity success"),
		}, nil
	}
	now := time.Now()
	cartEntity.Quantity = int(request.NewQuantity)
	cartEntity.UpdatedAt = &now
	cartEntity.UpdatedBy = &claims.FullName

	err = cs.cartRepository.UpdateCart(ctx, cartEntity)
	if err != nil {
		return nil, err
	}

	return &cart.UpdateCartQuantityResponse{
		Base: utils.SuccessResponse("Update cart quantity success"),
	}, nil

}

func NewCartService(productRepository repository.IProductRepository, cartRepository repository.ICartRepository) ICartService {
	return  &cartService{
		productRepository: productRepository,
		cartRepository: cartRepository,
	}
}
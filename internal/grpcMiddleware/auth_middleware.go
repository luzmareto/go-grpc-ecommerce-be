package grpcmiddleware

import (
	"context"
	"fmt"

	jwtentity "github.com/luzmareto/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

type authMiddleware struct {
	cacheService *gocache.Cache
}

func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	fmt.Println(info.FullMethod)
	if info.FullMethod == "/auth.AuthService/Login" || info.FullMethod == "/auth.AuthService/Register" || info.FullMethod == "/product.ProductService/DetailProduct" {
		return handler(ctx, req)

	}

	// ambil token dari metada
	tokenstr, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// cek token dari logout cache
	_, ok := am.cacheService.Get(tokenstr)
	if ok {
		return nil, utils.UnauthenticatedResponse()
	}

	// parse jwt hingga menjadi jwt
	claims, err := jwtentity.GetClaimsFromToken(tokenstr)
	if err != nil {
		return nil, err
	}

	// sematkan entity ke context
	ctx = claims.SetToContext(ctx)

	res, err := handler(ctx, req)

	return res, err
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}

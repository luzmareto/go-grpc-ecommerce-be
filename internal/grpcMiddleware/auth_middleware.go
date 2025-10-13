package grpcmiddleware

import (
	"context"

	jwtentity "github.com/luzmareto/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

type authMiddleware struct {
	cacheService *gocache.Cache
}

//api yang tidak perlu login
var publicApis = map[string]bool{
	 "/auth.AuthService/Login": true,
	 "/auth.AuthService/Register": true,
	 "/product.ProductService/DetailProduct": true,
	 "/product.ProductService/ListProduct": true,
	 "/product.ProductService/ListProducts": true,
	 "/product.ProductService/HighlightProducts": true,
	 "/product.ProductService/HighlightProduct": true,
}


func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if publicApis[info.FullMethod] {
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

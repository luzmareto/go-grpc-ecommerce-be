package main

// go run cmd/grpc/main.go

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv" //import manual
	grpcmiddleware "github.com/luzmareto/go-grpc-ecommerce-be/internal/grpcMiddleware"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/handler"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/repository"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/service"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/auth"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/cart"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/order"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/product"
	"github.com/luzmareto/go-grpc-ecommerce-be/pkg/database"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"            //import manual
	"google.golang.org/grpc/reflection" //import manual
)

func main() {
	ctx := context.Background()
	godotenv.Load()
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Panicf("error when listen %v", err)
	}

	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))
	log.Println("Connected to database")

	cacheService := gocache.New(time.Hour*24, time.Hour)

	authMiddleware := grpcmiddleware.NewAuthMiddleware(cacheService)

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, cacheService)
	authHandler := handler.NewAuthHandler(authService)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)
	
	cartRepository := repository.NewCartRepository(db)
	cartService := service.NewCartService(productRepository, cartRepository)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(db, orderRepository,productRepository)
	orderHandler := handler.NewOrderHandler(orderService)

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.ErrorMiddleware,
			authMiddleware.Middleware,
		),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)
	product.RegisterProductServiceServer(serv, productHandler)
	cart.RegisterCartServiceServer(serv, cartHandler)
	order.RegisterOrderServiceServer(serv, orderHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection is registered.")
	}

	log.Println("Server is runing on :50052 port.")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("Server is error %v", err)
	}
}

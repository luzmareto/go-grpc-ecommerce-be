package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv" //import manual
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/handler"
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/service"
	"github.com/luzmareto/go-grpc-ecommerce-be/pkg/database"
	"github.com/luzmareto/go-grpc-ecommerce-be/pkg/grpcmiddleware" //import manual
	"google.golang.org/grpc"                                       //import manual
	"google.golang.org/grpc/reflection"                            //import manual
)

func main() {
	ctx := context.Background()
	godotenv.Load()
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Panicf("error when listen %v", err)
	}

	database.ConnectDB(ctx, os.Getenv("DB_URI"))
	log.Println("Connected to database")

	serviceHandler := handler.NewServiceHandler()

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.ErrorMiddleware,
		),
	)

	service.RegisterHelloWorldServiceServer(serv, serviceHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection is registered.")
	}

	log.Println("Server is runing on :50052 port.")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("Server is error %v", err)
	}
}

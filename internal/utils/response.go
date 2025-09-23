package utils

import (
	"github.com/luzmareto/go-grpc-ecommerce-be/pb/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func BadRequestResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 400,
		Message:    message,
		IsError:    true,
	}
}

func NotFoundResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 404,
		Message:    message,
		IsError:    true,
	}
}

func UnauthenticatedResponse() error {
	return status.Error(codes.Unauthenticated, "Unauthenticated")
}

func ValidationErrorResponse(validationError []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:       400,
		Message:          "validation error",
		IsError:          true,
		ValidationErrors: validationError,
	}
}

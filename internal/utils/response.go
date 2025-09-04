package utils

import "github.com/luzmareto/go-grpc-ecommerce-be/pb/common"

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func ValidationErrorResponse(validationError []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:       400,
		Message:          "validation error",
		IsError:          true,
		ValidationErrors: validationError,
	}
}

package attribute

import (
	"calculationengine/constants"
	"calculationengine/models"
	"calculationengine/store"
	"context"

	"gorm.io/gorm"
)

func CreateAttribute(ctx context.Context, request models.CreateAttributeRequest) (*storage.ApiResponse, error) {
	result := gorm.WithResult()

	createObject := &storage.Attribute{Name: request.Name, DataType: request.DataType}

	err := gorm.G[storage.Attribute](storage.DB, result).Create(ctx, createObject)
	if err != nil {
		return &storage.ApiResponse{Message: "Something Went Wrong", Data: []any{}}, err
	}
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

func GetAllAttributes(ctx context.Context) (*models.GetAllAttributesResponse, error) {
	response := models.GetAllAttributesResponse{
		Message: constants.SUCCESS,
	}
	store := storage.NewStore(storage.DB)
	result, err := store.GetAllAttributes(ctx)
	if err != nil {
		response.Message = "Something went wrong"
		return &response, nil
	}
	response.Data = result
	return &response, nil
}

func GetCategoryWiseCommonAttributes(ctx context.Context, request models.GetCategoryWiseCommonAttributesRequest) (*models.GetCategoryWiseCommonAttributesResponse, error) {
	response := models.GetCategoryWiseCommonAttributesResponse{
		Message: constants.SUCCESS,
	}
	store := storage.NewStore(storage.DB)
	result, err := store.GetCategoryWiseCommonAttributes(ctx, request)
	if err != nil {
		response.Message = "Something went wrong"
		return &response, nil
	}
	response.Data = result
	return &response, nil
}

func ChangeCategoryAttributeAssignment(ctx context.Context, request models.ChangeCategoryAttributeAssignmentRequest) (*storage.ApiResponse, error) {
	store := storage.NewStore(storage.DB)
	err := store.ChangeCategoryAttributeAssignment(ctx, request)
	if err != nil {
		return &storage.ApiResponse{Message: "Something Went Wrong", Data: []any{}}, err
	}
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

func DeleteAttribute(ctx context.Context, request models.DeleteAttributeRequest) (*storage.ApiResponse, error) {
	s := storage.NewStore(storage.DB)
	if err := s.DeleteAttribute(ctx, request.AttributeID); err != nil {
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, err
	}
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

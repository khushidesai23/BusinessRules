package category

import (
	"calculationengine/models"
	storage "calculationengine/store"
	"context"
	"fmt"
)

func CreateCategory(ctx context.Context, request models.CreateCategoryRequest) (*storage.ApiResponse, error) {
	s := storage.NewStore(storage.DB)
	fmt.Println("Hello1")
	err := s.CreateCategory(ctx, request.Name)
	if err != nil {
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, nil
	}
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

func GetAllCategories(ctx context.Context) (*models.GetCategoriesResponse, error) {
	response := models.GetCategoriesResponse{}
	s := storage.NewStore(storage.DB)
	data, err := s.GetAllCategories(ctx)
	if err != nil {
		response.Message = "Something went wrong"
		return &response, nil
	}
	response.Message = "success"
	response.Data = data
	return &response, nil
}

func DeleteCategory(ctx context.Context, request models.DeleteCategoryRequest) (*storage.ApiResponse, error) {
	s := storage.NewStore(storage.DB)
	if err := s.DeleteCategory(ctx, request.CategoryID); err != nil {
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, err
	}
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

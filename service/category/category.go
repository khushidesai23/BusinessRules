package category

import (
	"calculationengine/models"
	storage "calculationengine/store"
	"context"
	"fmt"
)

// Fixed: CreateCategory now prevents FK violations by not auto-assigning default attributes
// Previously: Attempted to create CategoryAttributeAssignment with AttributeID=0 (which doesn't exist as a real Attribute)
// Now: Only creates the Category, allowing admin to explicitly assign real attributes later
func CreateCategory(ctx context.Context, request models.CreateCategoryRequest) (*storage.ApiResponse, error) {
	s := storage.NewStore(storage.DB)
	fmt.Println("Hello1")
	err := s.CreateCategory(ctx, request.Name)
	if err != nil {
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, nil
	}
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

// GetAllCategories retrieves all available categories
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

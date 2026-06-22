package category

import (
	"calculationengine/logging"
	"calculationengine/models"
	storage "calculationengine/store"
	"context"
)

// Fixed: CreateCategory now prevents FK violations by not auto-assigning default attributes
// Previously: Attempted to create CategoryAttributeAssignment with AttributeID=0 (which doesn't exist as a real Attribute)
// Now: Only creates the Category, allowing admin to explicitly assign real attributes later
func CreateCategory(ctx context.Context, request models.CreateCategoryRequest) (*storage.ApiResponse, error) {
	logger := logging.FromContext(ctx)
	s := storage.NewStore(storage.DB)
	logger.InfoContext(ctx, "create category requested", "category_name", request.Name)
	err := s.CreateCategory(ctx, request.Name)
	if err != nil {
		logger.ErrorContext(ctx, "create category failed", "category_name", request.Name, "error", err)
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, nil
	}
	logger.InfoContext(ctx, "category created", "category_name", request.Name)
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

// GetAllCategories retrieves all available categories
func GetAllCategories(ctx context.Context) (*models.GetCategoriesResponse, error) {
	response := models.GetCategoriesResponse{}
	s := storage.NewStore(storage.DB)
	data, err := s.GetAllCategories(ctx)
	if err != nil {
		logging.FromContext(ctx).ErrorContext(ctx, "get categories failed", "error", err)
		response.Message = "Something went wrong"
		return &response, nil
	}
	logging.FromContext(ctx).InfoContext(ctx, "categories fetched", "count", len(data))
	response.Message = "success"
	response.Data = data
	return &response, nil
}

func DeleteCategory(ctx context.Context, request models.DeleteCategoryRequest) (*storage.ApiResponse, error) {
	s := storage.NewStore(storage.DB)
	if err := s.DeleteCategory(ctx, request.CategoryID); err != nil {
		logging.FromContext(ctx).ErrorContext(ctx, "delete category failed", "category_id", request.CategoryID, "error", err)
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, err
	}
	logging.FromContext(ctx).InfoContext(ctx, "category deleted", "category_id", request.CategoryID)
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

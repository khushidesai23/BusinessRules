package attribute

import (
	"calculationengine/constants"
	"calculationengine/logging"
	"calculationengine/models"
	"calculationengine/store"
	"context"

	"gorm.io/gorm"
)

func CreateAttribute(ctx context.Context, request models.CreateAttributeRequest) (*storage.ApiResponse, error) {
	logger := logging.FromContext(ctx)
	result := gorm.WithResult()

	createObject := &storage.Attribute{Name: request.Name, DataType: request.DataType}
	logger.InfoContext(ctx, "create attribute requested", "attribute_name", request.Name, "data_type", request.DataType)

	err := gorm.G[storage.Attribute](storage.DB, result).Create(ctx, createObject)
	if err != nil {
		logger.ErrorContext(ctx, "create attribute failed", "attribute_name", request.Name, "data_type", request.DataType, "error", err)
		return &storage.ApiResponse{Message: "Something Went Wrong", Data: []any{}}, err
	}
	logger.InfoContext(ctx, "attribute created", "attribute_id", createObject.ID, "attribute_name", request.Name, "data_type", request.DataType)
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

func GetAllAttributes(ctx context.Context) (*models.GetAllAttributesResponse, error) {
	response := models.GetAllAttributesResponse{
		Message: constants.SUCCESS,
	}
	store := storage.NewStore(storage.DB)
	result, err := store.GetAllAttributes(ctx)
	if err != nil {
		logging.FromContext(ctx).ErrorContext(ctx, "get attributes failed", "error", err)
		response.Message = "Something went wrong"
		return &response, nil
	}
	logging.FromContext(ctx).InfoContext(ctx, "attributes fetched", "count", len(result))
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
		logging.FromContext(ctx).ErrorContext(ctx, "get common attributes failed", "category_ids", request.CategoryIDs, "error", err)
		response.Message = "Something went wrong"
		return &response, nil
	}
	logging.FromContext(ctx).InfoContext(ctx, "common attributes fetched", "category_ids", request.CategoryIDs, "count", len(result))
	response.Data = result
	return &response, nil
}

func ChangeCategoryAttributeAssignment(ctx context.Context, request models.ChangeCategoryAttributeAssignmentRequest) (*storage.ApiResponse, error) {
	store := storage.NewStore(storage.DB)
	logging.FromContext(ctx).InfoContext(ctx, "change category attribute assignment requested",
		"assign_category_ids", request.Assign.CategoryIDs,
		"assign_attribute_ids", request.Assign.AttributeIDs,
		"unassign_category_ids", request.UnAssign.CategoryIDs,
		"unassign_attribute_ids", request.UnAssign.AttributeIDs,
	)
	err := store.ChangeCategoryAttributeAssignment(ctx, request)
	if err != nil {
		logging.FromContext(ctx).ErrorContext(ctx, "change category attribute assignment failed", "error", err)
		return &storage.ApiResponse{Message: "Something Went Wrong", Data: []any{}}, err
	}
	logging.FromContext(ctx).InfoContext(ctx, "category attribute assignment changed")
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

func DeleteAttribute(ctx context.Context, request models.DeleteAttributeRequest) (*storage.ApiResponse, error) {
	s := storage.NewStore(storage.DB)
	if err := s.DeleteAttribute(ctx, request.AttributeID); err != nil {
		logging.FromContext(ctx).ErrorContext(ctx, "delete attribute failed", "attribute_id", request.AttributeID, "error", err)
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, err
	}
	logging.FromContext(ctx).InfoContext(ctx, "attribute deleted", "attribute_id", request.AttributeID)
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

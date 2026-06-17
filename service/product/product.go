package product

import (
	"calculationengine/logging"
	"calculationengine/models"
	"calculationengine/service/formulas"
	"calculationengine/service/utils"
	storage "calculationengine/store"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func UpsertProduct(ctx context.Context, request models.CreateProductRequest) (*storage.ApiResponse, error) {
	logger := logging.FromContext(ctx)
	s := storage.NewStore(storage.DB)
	attributes, err := s.GetAttributesIdDataMap(ctx, []int{request.CategryID})
	var id string
	create := false
	if request.ProductID == "" {
		create = true
		id = uuid.NewString()
	} else {
		id = request.ProductID
	}
	if err != nil {
		logger.ErrorContext(ctx, "load category attributes failed", "category_id", request.CategryID, "error", err)
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, nil
	}
	logger.InfoContext(ctx, "product upsert requested",
		"category_id", request.CategryID,
		"product_id", id,
		"create", create,
		"payload_count", len(request.ProductData),
	)
	var createProductParams []models.CreateProductParams
	defaultAttributePresent := false

	// Enhanced: Process product data with support for attributeId 0 (product name)
	for _, data := range request.ProductData {
		value := strings.TrimSpace(data.Value)

		if data.AttributeID == 0 {
			// Fixed: AttributeID 0 is the special "product name" field
			// Skip lookup in attributes map since it's not a real Attribute record
			// Product names don't require FK validation
			if value != "" {
				defaultAttributePresent = true
			}
			// Allow saving the product name value without type validation
		} else {
			// For actual attributes (ID > 0), validate against the category's attributes
			attribute, ok := attributes[data.AttributeID]
			if !ok {
				logger.WarnContext(ctx, "product upsert rejected because attribute is not assigned to category", "category_id", request.CategryID, "attribute_id", data.AttributeID)
				return &storage.ApiResponse{Message: fmt.Sprintf("Attribute %d does not exist in this category", data.AttributeID), Data: []any{}}, nil
			}
			validationErr := validateData(value, attribute)
			if validationErr != nil {
				logger.WarnContext(ctx, "product upsert rejected because attribute validation failed", "attribute_id", data.AttributeID, "data_type", attribute.DataType, "error", validationErr)
				return &storage.ApiResponse{Message: validationErr.Error(), Data: []any{}}, nil
			}
		}
		createProductParams = append(createProductParams, models.CreateProductParams{
			ID:          id,
			CategoryID:  uint(request.CategryID),
			AttributeID: uint(data.AttributeID),
			Data:        value,
		})
	}

	// Fixed: Require product name for new product creation
	// Product name (attributeId=0) is mandatory to distinguish products
	if create && !defaultAttributePresent {
		logger.WarnContext(ctx, "product create rejected because product name is missing", "category_id", request.CategryID)
		return &storage.ApiResponse{Message: "Please enter Product Name to create a new product", Data: []any{}}, nil
	}
	if err := s.UpsertProduct(ctx, createProductParams); err != nil {
		logger.ErrorContext(ctx, "product upsert failed", "product_id", id, "category_id", request.CategryID, "error", err)
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, err
	}
	logger.InfoContext(ctx, "product values stored", "product_id", id, "records", len(createProductParams))

	// After saving product data, evaluate formulas to compute derived attributes
	evaluateFormulaRequest := models.EvaluateFormulaRequest{
		ProductID: []string{id},
	}
	evaluatedProductData, err := formulas.EvaluateFormula(ctx, evaluateFormulaRequest)
	if err != nil {
		logger.ErrorContext(ctx, "product formula evaluation failed", "product_id", id, "error", err)
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, err
	}
	logger.InfoContext(ctx, "product formulas evaluated", "product_id", id, "computed_records", len(evaluatedProductData))
	// Save computed formula results back to product
	if err := s.UpsertProduct(ctx, evaluatedProductData); err != nil {
		logger.ErrorContext(ctx, "persist computed product values failed", "product_id", id, "error", err)
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, err
	}
	logger.InfoContext(ctx, "product upsert completed", "product_id", id, "category_id", request.CategryID)
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

// Enhanced: Retrieve all products in a list format with product names and category info
// Product names are now stored in the products.name column (attributeId=0)
func GetProductList(ctx context.Context) (*models.GetProductListResponse, error) {
	response := models.GetProductListResponse{}
	s := storage.NewStore(storage.DB)
	data, err := s.GetProductList(ctx)
	if err != nil {
		logging.FromContext(ctx).ErrorContext(ctx, "get products failed", "error", err)
		response.Message = "Something went wrong"
		return &response, nil
	}
	logging.FromContext(ctx).InfoContext(ctx, "products fetched", "count", len(data))
	response.Message = "success"
	response.Data = data
	return &response, nil
}

// Enhanced: Retrieve single product with all attribute data including product name
// Product name is returned as attributeId=0 in the result set
func GetSingleProductData(ctx context.Context, request models.GetProductDataRequest) (*models.GetProductDataResponse, error) {
	response := models.GetProductDataResponse{}
	s := storage.NewStore(storage.DB)
	data, err := s.GetProductData(ctx, []string{request.ProductID})
	if err != nil {
		logging.FromContext(ctx).ErrorContext(ctx, "get product data failed", "product_id", request.ProductID, "error", err)
		response.Message = "Something went wrong"
		return &response, nil
	}
	logging.FromContext(ctx).InfoContext(ctx, "product data fetched", "product_id", request.ProductID, "records", len(data))
	response.Message = "success"
	response.Data = data
	return &response, nil

}

func validateData(value string, attribute storage.Attribute) error {
	if value == "" {
		return nil
	}
	switch attribute.DataType {
	case "integer":
		_, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
	case "float":
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
	case "boolean":
		_, err := utils.StringToBoolean(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteProduct(ctx context.Context, request models.DeleteProductRequest) (*storage.ApiResponse, error) {
	s := storage.NewStore(storage.DB)
	if err := s.DeleteProduct(ctx, request.ProductID); err != nil {
		logging.FromContext(ctx).ErrorContext(ctx, "delete product failed", "product_id", request.ProductID, "error", err)
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, err
	}
	logging.FromContext(ctx).InfoContext(ctx, "product deleted", "product_id", request.ProductID)
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

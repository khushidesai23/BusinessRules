package product

import (
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
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, nil
	}
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
				return &storage.ApiResponse{Message: fmt.Sprintf("Attribute %d does not exist in this category", data.AttributeID), Data: []any{}}, nil
			}
			validationErr := validateData(value, attribute)
			if validationErr != nil {
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
		return &storage.ApiResponse{Message: "Please enter Product Name to create a new product", Data: []any{}}, nil
	}
	s.UpsertProduct(ctx, createProductParams)
	
	// After saving product data, evaluate formulas to compute derived attributes
	evaluateFormulaRequest := models.EvaluateFormulaRequest{
		ProductID: []string{id},
	}
	evaluatedProductData, _ := formulas.EvaluateFormula(ctx, evaluateFormulaRequest)
	fmt.Println(evaluatedProductData)
	// Save computed formula results back to product
	s.UpsertProduct(ctx, evaluatedProductData)
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

// Enhanced: Retrieve all products in a list format with product names and category info
// Product names are now stored in the products.name column (attributeId=0)
func GetProductList(ctx context.Context) (*models.GetProductListResponse, error) {
	response := models.GetProductListResponse{}
	s := storage.NewStore(storage.DB)
	data, err := s.GetProductList(ctx)
	if err != nil {
		response.Message = "Something went wrong"
		return &response, nil
	}
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
		response.Message = "Something went wrong"
		return &response, nil
	}
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

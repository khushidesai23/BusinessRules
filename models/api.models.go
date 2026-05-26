package models

type CreateAttributeRequest struct {
	Name     string `json:"name" validate:"required"`
	DataType string `json:"dataType" validate:"required"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type GetCategoriesResponse struct {
	Message string                `json:"message"`
	Data    []GetCategoriesResult `json:"data"`
}

type GetCategoryWiseCommonAttributesRequest struct {
	CategoryIDs []int `json:"categoryIds" validate:"required"`
}

type GetCategoryWiseCommonAttributesResponse struct {
	Message string                                  `json:"message"`
	Data    []GetCategoryWiseCommonAttributesResult `json:"data"`
}

type GetAllAttributesResponse struct {
	Message string                `json:"message"`
	Data    []GetAttributesResult `json:"data"`
}

type GetAllFormulasResponse struct {
	Message string               `json:"message"`
	Data    []FormulasListResult `json:"data"`
}

type GetProductDataRequest struct {
	ProductID string `json:"productId"`
}

type GetProductDataResponse struct {
	Message string               `json:"message"`
	Data    []ProductDatasResult `json:"data"`
}

type GetProductListResponse struct {
	Message string              `json:"message"`
	Data    []ProductListResult `json:"data"`
}

type ChangeCategoryAttributeAssignmentRequest struct {
	Assign struct {
		CategoryIDs  []int `json:"categoryIds"`
		AttributeIDs []int `json:"attributeIds"`
	} `json:"assign"`
	UnAssign struct {
		CategoryIDs  []int `json:"categoryIds"`
		AttributeIDs []int `json:"attributeIds"`
	} `json:"unassign"`
}

type CreateFormulaRequest struct {
	CategoryID      int    `json:"categoryId" validate:"required"`
	TargetAttribute int    `json:"targetAttribute" validate:"required"`
	Formula         string `json:"formula" validate:"required"`
}

type EvaluateFormulaRequest struct {
	ProductID []string `json:"poductIds" validate:"required"`
}

type CreateProductRequest struct {
	CategryID   int    `json:"categoryId" validate:"required"`
	ProductID   string `json:"productId"`
	ProductData []struct {
		AttributeID int    `json:"attributeId" validate:"required"`
		Value       string `json:"value" validate:"required"`
	} `json:"data"`
}

type DeleteCategoryRequest struct {
	CategoryID int `json:"categoryId" validate:"required"`
}

type DeleteAttributeRequest struct {
	AttributeID int `json:"attributeId" validate:"required"`
}

type DeleteProductRequest struct {
	ProductID string `json:"productId" validate:"required"`
}

type DeleteFormulaRequest struct {
	CategoryID      int `json:"categoryId" validate:"required"`
	TargetAttribute int `json:"targetAttribute" validate:"required"`
}

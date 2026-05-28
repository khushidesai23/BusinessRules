package models

type DataType string

const (
	Integer DataType = "integer"
	String  DataType = "string"
	Boolean DataType = "boolean"
	Float   DataType = "float"
)

type AttributeDependenciesResult struct {
	TargetAttributeID    int `json:"targetAttributeId"`
	DependentAttributeID int `json:"dependentAttributeId"`
	CategoryID           int `json:"categoryId"`
}

type FormulasListResult struct {
	CategoryID          int    `json:"categoryId"`
	TargetAttributeID   int    `json:"targetAttributeId"`
	CategoryName        string `json:"categoryName"`
	TargetAttributeName string `json:"targetAttributeName"`
	Formula             string `json:"formula"`
}

type ProductDatasResult struct {
	CategoryID    int      `json:"categoryId"`
	ID            string   `json:"id"`
	AttributeID   int      `json:"attributeId"`
	// Enhanced: Can be product name (when attributeId=0) or attribute data value
	Data          string   `json:"data"`
	DataType      DataType `json:"dataType"`
	AttributeName string   `json:"attributeName"`
}

// Enhanced: ProductListResult now includes product Name field for display purposes
// The Name field is populated from products.name column (stored for attributeId=0)
type ProductListResult struct {
	ID           string `json:"id"`
	// New: Product name (stored separately as attributeId=0 in products table)
	Name         string `json:"name"`
	CategoryID   int    `json:"categoryId"`
	CategoryPath string `json:"categoryPath"`
}

type GetCategoriesResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetCategoryWiseCommonAttributesResult struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
	Assigned bool   `json:"assigned"`
}

type GetAttributesResult struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type FormulasResult struct {
	CategoryID        int    `json:"categoryId"`
	TargetAttributeID int    `json:"targetAttributeId"`
	Expression        string `json:"expression"`
}

type SaveFormulaParams struct {
	CategoryID                      int
	TopologicallySortedAttributeIDs []int
	Formula                         string
	DependentAttributeIDs           []int
	TargetAttributeID               int
}

// Enhanced: CreateProductParams now supports both product names and attribute values
// - When AttributeID=0: Data contains the product name (stored in products.name column)
// - When AttributeID>0: Data contains attribute value (stored in products.data column)
// This unified structure allows flexible upsertion of both product names and attribute data
type CreateProductParams struct {
	ID          string
	CategoryID  uint
	AttributeID uint // 0 for product name, >0 for attribute values
	Data        string
}

type TopologicalSortResult struct {
	CategoryID           int `json:"categoryId"`
	AttributeID          int `json:"attributeId"`
	TopologicalSortOrder int `json:"topologicalSortOrder"`
}

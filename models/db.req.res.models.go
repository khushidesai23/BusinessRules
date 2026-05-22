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
	Data          string   `json:"data"`
	DataType      DataType `json:"dataType"`
	AttributeName string   `json:"attributeName"`
}

type ProductListResult struct {
	ID           string `json:"id"`
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

type CreateProductParams struct {
	ID          string
	CategoryID  uint
	AttributeID uint
	Data        string
}

type TopologicalSortResult struct {
	CategoryID           int `json:"categoryId"`
	AttributeID          int `json:"attributeId"`
	TopologicalSortOrder int `json:"topologicalSortOrder"`
}

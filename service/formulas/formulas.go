package formulas

import (
	"calculationengine/models"
	"calculationengine/service/dag"
	"calculationengine/service/evaluator"
	"calculationengine/service/parser"
	"calculationengine/service/utils"
	storage "calculationengine/store"
	"context"
	"errors"
	"strconv"

	// "errors"
	"fmt"
	"slices"
	// "gorm.io/gorm"
)

type Formula struct {
	categoryId int
	formula    string
	attributes []string
}

func CreateFormula(ctx context.Context, request models.CreateFormulaRequest) (*storage.ApiResponse, error) {
	attributesInFormula, err := getAttributesFromFormula(request.Formula)
	if err != nil {
		return &storage.ApiResponse{Message: err.Error(), Data: []any{}}, nil
	}
	attributesInFormula = utils.RemoveArrayDuplicates(attributesInFormula)

	var store = storage.NewStore(storage.DB)
	attributesAssignedToCategory := store.GetAttributesByNames(ctx, request.CategoryID, attributesInFormula)
	attributeNameAssignedToCategory := utils.Map(attributesAssignedToCategory, func(a storage.Attribute) string {
		return a.Name
	})

	attributesNotExisting := utils.ArrayDifference(attributesInFormula, attributeNameAssignedToCategory)

	if len(attributesNotExisting) != 0 {
		return &storage.ApiResponse{Message: fmt.Sprintf("%s Attributes does not exist", attributesNotExisting), Data: []any{}}, nil
		//!Handle error
	}
	attributeDependencies := store.GetAllFormulaDependencies(ctx, []int{request.CategoryID})

	//^Remove duplicates from this array
	var attributeIdsInFormula []int
	for _, attribute := range attributesAssignedToCategory {
		if slices.Contains(attributesInFormula, attribute.Name) {
			attributeIdsInFormula = append(attributeIdsInFormula, int(attribute.ID))
		}
	}

	for _, dependentAttributeId := range attributeIdsInFormula {
		isDependencyExisting := false
		for _, attributeDependency := range attributeDependencies {
			if attributeDependency.TargetAttributeID == request.TargetAttribute && attributeDependency.DependentAttributeID == dependentAttributeId {
				isDependencyExisting = true
				break
			}
		}
		if isDependencyExisting {
			continue
		}
		attributeDependencies = append(attributeDependencies, models.AttributeDependenciesResult{TargetAttributeID: request.TargetAttribute, DependentAttributeID: dependentAttributeId})
	}

	graph := generateGraphFromDependencies(attributeDependencies)
	topologicalSortedAttributes, hasCycle := graph.TopologicalSort()
	if hasCycle {
		fmt.Println("Has Cycle")
		return &storage.ApiResponse{Message: "Cycle Detected", Data: []any{}}, nil
	}

	_, err2 := checkFormulaSyntaxErrors(request.Formula, attributesAssignedToCategory)
	if err2 != nil {
		return &storage.ApiResponse{Message: err2.Error(), Data: []any{}}, nil
	}

	err1 := store.SaveFormula(ctx, models.SaveFormulaParams{
		CategoryID:                      request.CategoryID,
		TopologicallySortedAttributeIDs: topologicalSortedAttributes,
		Formula:                         request.Formula,
		DependentAttributeIDs:           attributeIdsInFormula,
		TargetAttributeID:               request.TargetAttribute,
	})
	if err1 != nil {
		return &storage.ApiResponse{Message: "Something Went wrong", Data: []any{}}, nil
	}
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

func GetFormulasList(ctx context.Context) (*models.GetAllFormulasResponse, error) {
	response := models.GetAllFormulasResponse{}
	s := storage.NewStore(storage.DB)
	data, err := s.GetFormulasList(ctx)
	if err != nil {
		response.Message = "Something went wrong"
		return &response, nil
	}
	response.Message = "success"
	response.Data = data
	return &response, nil
}

func checkFormulaSyntaxErrors(formula string, attributes []storage.Attribute) (bool, error) {
	//Generating Default Env
	env := evaluator.NewEnvironment()
	for _, attribute := range attributes {
		switch attribute.DataType {
		case "integer":
			object := &evaluator.Integer{Value: 1}
			env.Set(attribute.Name, object)
		case "boolean":
			object := evaluator.NativeBoolToBooleanObject(true)
			env.Set(attribute.Name, object)
		case "string":
			object := &evaluator.String{Value: "Hello"}
			env.Set(attribute.Name, object)
		case "float":
			object := &evaluator.Float{Value: 1.5}
			env.Set(attribute.Name, object)
		}
	}
	result := parseAndEvaluateFormula(formula, env)
	if result.Type() == evaluator.ERROR_OBJ {
		return false, errors.New(result.Inspect())
	}
	return true, nil
}

func EvaluateFormula(ctx context.Context, request models.EvaluateFormulaRequest) ([]models.CreateProductParams, error) {
	response := []models.CreateProductParams{}
	var store = storage.NewStore(storage.DB)
	productDatas, err := store.GetProductData(ctx, request.ProductID)
	if err != nil {
		return response, nil
	}
	categoryIds := utils.Map(productDatas, func(productData models.ProductDatasResult) int {
		return productData.CategoryID
	})
	categoryIds = utils.RemoveArrayDuplicates(categoryIds)
	formulaDependencies := store.GetAllFormulaDependencies(ctx, categoryIds)
	attributesIdMap, err := store.GetAttributesIdDataMap(ctx, categoryIds)
	if err != nil {

	}
	formulas := store.GetFormulas(ctx, categoryIds)
	topologicalSortOrder, err := store.GetTopologicalSorting(ctx, categoryIds)
	if err != nil {
		return response, nil
	}

	for _, productId := range request.ProductID {
		productData := utils.Filter(productDatas, func(productData models.ProductDatasResult) bool {
			return productData.ID == productId
		})
		if len(productData) == 0 {
			continue
		}
		productCategoryId := productData[0].CategoryID
		env := generateEnvironmentFromProductData(productData)
		categoryTopologicalSortOrder := utils.Filter(topologicalSortOrder, func(sortOrder models.TopologicalSortResult) bool {
			return sortOrder.CategoryID == productCategoryId
		})
		for _, attribute := range categoryTopologicalSortOrder {
			attributeId := attribute.AttributeID
			formula := utils.Filter(formulas, func(x models.FormulasResult) bool {
				return x.CategoryID == productCategoryId && x.TargetAttributeID == attributeId
			})
			if len(formula) == 0 {
				continue
			}
			currentFormulaDependencies := utils.Filter(formulaDependencies, func(x models.AttributeDependenciesResult) bool {
				return x.CategoryID == productCategoryId && x.TargetAttributeID == formula[0].TargetAttributeID
			})
			nilValuePresent := false
			for _, currentFormulaDependency := range currentFormulaDependencies {
				dependentAttributeName := attributesIdMap[currentFormulaDependency.DependentAttributeID].Name
				if _, ok := env.Get(dependentAttributeName); !ok {
					nilValuePresent = true
					break
				}
			}
			if nilValuePresent {
				continue
			}
			obj := parseAndEvaluateFormula(formula[0].Expression, env)
			if obj.Type() == evaluator.ERROR_OBJ {
				continue
			}
			response = append(response, models.CreateProductParams{
				ID:          productId,
				AttributeID: uint(attributeId),
				Data:        obj.Inspect(),
				CategoryID:  uint(productCategoryId),
			})
			targetAttributeName := attributesIdMap[attributeId].Name
			env.Set(targetAttributeName, obj)
		}

	}
	return response, nil
}

func parseAndEvaluateFormula(formula string, env *evaluator.Environment) evaluator.Object {
	lexer := parser.NewLexer(formula)
	nparser := parser.NewParser(lexer)
	program := nparser.ParseProgram()
	eval := evaluator.Eval(program.Statements[0], env)
	return eval
}

func generateEnvironmentFromProductData(productData []models.ProductDatasResult) *evaluator.Environment {
	env := evaluator.NewEnvironment()
	for _, data := range productData {
		if data.Data == "" {
			continue
		}
		switch data.DataType {
		case "integer":
			value, err := strconv.Atoi(data.Data)
			if err != nil {
				continue
			}
			object := &evaluator.Integer{Value: int64(value)}
			env.Set(data.AttributeName, object)
		case "boolean":
			value, err := utils.StringToBoolean(data.Data)
			if err != nil {
				continue
			}
			object := evaluator.NativeBoolToBooleanObject(value)
			env.Set(data.AttributeName, object)
		case "string":
			object := &evaluator.String{Value: string(data.Data)}
			fmt.Println("Object: ", object)
			env.Set(data.AttributeName, object)
		case "float":
			value, err := strconv.ParseFloat(data.Data, 64)
			if err != nil {
				continue
			}
			object := &evaluator.Float{Value: float64(value)}
			env.Set(data.AttributeName, object)
		}
	}
	return env
}

func generateGraphFromDependencies(attributeDependencies []models.AttributeDependenciesResult) *dag.GraphList {
	graph := dag.NewGraphList()
	for _, value := range attributeDependencies {
		graph.AddEdge(value.DependentAttributeID, value.TargetAttributeID)
	}
	return graph
}

func getAttributesFromFormula(formula string) ([]string, error) {
	var attributesInFormula []string
	lexer := parser.NewLexer(formula)
	illegalTokens := []string{}
formulaTokenLoop:
	for {
		var token = lexer.NextToken()
		switch token.TokenType {
		case parser.EOF:
			break formulaTokenLoop
		case parser.ILLEGAL:
			illegalTokens = append(illegalTokens, token.TokenValue)
			break formulaTokenLoop
			//!Handle error
		case parser.IDENT:
			attributesInFormula = append(attributesInFormula, token.TokenValue)
		}
	}
	if len(illegalTokens) != 0 {
		return []string{}, fmt.Errorf("%s nt allowed in formula", illegalTokens)
	}
	return attributesInFormula, nil
}

func (formula *Formula) validateAttributes() {

}

func DeleteFormula(ctx context.Context, request models.DeleteFormulaRequest) (*storage.ApiResponse, error) {
	s := storage.NewStore(storage.DB)
	if err := s.DeleteFormula(ctx, request.CategoryID, request.TargetAttribute); err != nil {
		return &storage.ApiResponse{Message: "Something went wrong", Data: []any{}}, err
	}
	return &storage.ApiResponse{Message: "success", Data: []any{}}, nil
}

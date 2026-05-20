package storage

import (
	"calculationengine/models"
	"context"
	"fmt"
	"strings"

	// "fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	// "gorm.io/gorm/clause"
)

type Store struct {
	DB *gorm.DB
}

func NewStore(db *gorm.DB) *Store{
	return &Store{DB:db}
}

func (s *Store) CreateCategory(ctx context.Context, name string) error{
	fmt.Println("Hello2")
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		category := Category{
			Path:name,
		}
		fmt.Println("Hello3")
		err := tx.Save(&category).Error
		fmt.Println("err: ", err)
		if err!=nil{
			return err
		}
		// Do not create a CategoryAttributeAssignment with AttributeID 0
		// (was causing FK violations because attribute id 0 does not exist).
		return nil
	})
	return err
}

func (s *Store) GetAllCategories(ctx context.Context) ([]models.GetCategoriesResult, error) {
	var categories []models.GetCategoriesResult
	err := s.DB.Model(&Category{}).
		Select("id, path as name").
		Order("updated_at DESC").
		Scan(&categories).Error
	return categories, err
}

func (s *Store) GetAllAttributes(ctx context.Context) ([]models.GetAttributesResult, error) {
	var attributes []models.GetAttributesResult
	err := s.DB.Model(&Attribute{}).
		Select(`id, name, data_type as "dataType"`).
		Order("updated_at DESC").
		Scan(&attributes).Error
	return attributes, err
}

func (s *Store) GetCategoryWiseCommonAttributes(ctx context.Context, params models.GetCategoryWiseCommonAttributesRequest) ([]models.GetCategoryWiseCommonAttributesResult, error) {
	var attributes []models.GetCategoryWiseCommonAttributesResult
	err := s.DB.Raw(`
		SELECT 
			a.id,
			a.name,
			a.data_type as "dataType",
			CASE 
				WHEN COUNT(DISTINCT caa.category_id) = ? THEN TRUE -- Replace 3 with the length of your array
				ELSE FALSE 
			END AS assigned
		FROM 
			attributes a
		LEFT JOIN 
			category_attribute_assignments caa 
			ON a.id = caa.attribute_id 
			AND caa.category_id IN ? -- Replace with your array elements
		GROUP BY 
			a.id, 
			a.name;
	`, len(params.CategoryIDs), params.CategoryIDs).Scan(&attributes).Error
	return attributes, err
}

func (s *Store) ChangeCategoryAttributeAssignment(ctx context.Context, params models.ChangeCategoryAttributeAssignmentRequest) error {
	var assignments []CategoryAttributeAssignment
	for _, catID := range params.Assign.CategoryIDs {
		for _, attrID := range params.Assign.AttributeIDs {
			assignments = append(assignments, CategoryAttributeAssignment{
				CategoryID: uint(catID),
				AttributeID: uint(attrID),
			})
		}
	}
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if len(assignments) != 0 {
			err := tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&assignments).Error
			if err != nil {
				return err
			}
		}
		if len(params.UnAssign.CategoryIDs) != 0 && len(params.UnAssign.AttributeIDs) != 0 {
			err := tx.Unscoped().
				Where("category_id IN ? AND attribute_id IN ?", params.UnAssign.CategoryIDs, params.UnAssign.AttributeIDs).
				Delete(&CategoryAttributeAssignment{}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *Store) UpsertProduct(ctx context.Context, datas []models.CreateProductParams) error {
	if len(datas) == 0 {
		return nil
	}
	var insertValues []string
	for _, data := range datas {
		insertValues = append(insertValues, fmt.Sprintf(`('%s', %d, %d, '%s')`, data.ID, data.CategoryID, data.AttributeID, data.Data)) 
	}
	queryStr := strings.Join(insertValues, ",")
	query := fmt.Sprintf(`
		INSERT INTO products (id, category_id, attribute_id, data)
		VALUES
			%s
		ON CONFLICT (id, category_id, attribute_id)
		DO UPDATE SET
			data = EXCLUDED.data
	`, queryStr)
	fmt.Println(query)
	err := s.DB.Exec(query).Error
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAttributesByNames(ctx context.Context, categoryId int, names []string) []Attribute {
	var attributes []Attribute
	s.DB.Model(&CategoryAttributeAssignment{}).
		Select(`attributes.id as "id", attributes.name as "name", attributes.data_type as "dataType"`).
		Joins("JOIN attributes ON attributes.id = category_attribute_assignments.attribute_id AND category_attribute_assignments.category_id = ? AND attributes.name in ?", categoryId, names).
		Scan(&attributes)
	return attributes
}

func (s *Store) GetAllFormulaDependencies(ctx context.Context, categoryIds []int) []models.AttributeDependenciesResult {
	var formulaDependencies []models.AttributeDependenciesResult
	s.DB.Model(&FormulaDependencies{}).
		Select(`category_id as "categoryId", target_attribute_id as "targetAttributeId", dependent_attribute_id as "dependentAttributeId"`).
		Where("category_id IN ?", categoryIds).
		Scan(&formulaDependencies)
	return formulaDependencies
}

func (s *Store) GetAttributesIdDataMap(ctx context.Context, categoryIds []int) (map[int]Attribute, error){
	var attributes []Attribute
	err := s.DB.Raw(`
		select 
			id, 
			name, 
			data_type as "dataType" 
		from 
			attributes 
		where 
			id in (
				select 
					attribute_id 
				from 
					category_attribute_assignments 
				where 
					category_id in ?
			)
	`, categoryIds).Scan(&attributes).Error
	attributesMap := make(map[int]Attribute, len(attributes))
	if err != nil {
		return attributesMap, err
	}
	for _, attribute := range attributes {
		attributesMap[int(attribute.ID)] = attribute
	}
	return attributesMap, nil
}

func (s *Store) GetAllFormulas(ctx context.Context, categoryIds []int) ([]models.FormulasResult, error) {
	var formulas []models.FormulasResult
	err := s.DB.Model(&Formulas{}).
		Select(`category_id as "categoryId", target_attribute_id as "targetAttributeId", expression`).
		Where("category_id IN ?", categoryIds).
		Scan(&formulas).Error
	return formulas, err
}

func (s *Store) SaveFormula(ctx context.Context, params models.SaveFormulaParams) error {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		gorm.G[CategoryAttributeAssignment](tx).Where("category_id = ?", params.CategoryID).Update(ctx, "topological_sort_order", nil)
		for index, value := range params.TopologicallySortedAttributeIDs{
			gorm.G[CategoryAttributeAssignment](tx).Where("category_id = ? AND attribute_id = ?", params.CategoryID, value).Update(ctx, "topological_sort_order", index)
		}
		tx.Save(&Formulas{
			CategoryID: uint(params.CategoryID),
			Expression: params.Formula,
			TargetAttributeID: uint(params.TargetAttributeID),
		})
		var formulaDependencies []FormulaDependencies
		for _, dependentAttributeId := range params.DependentAttributeIDs{
			formulaDependency := FormulaDependencies{
				CategoryID: uint(params.CategoryID),
				TargetAttributeID: uint(params.TargetAttributeID),
				DependentAttributeID: uint(dependentAttributeId),
			}
			formulaDependencies = append(formulaDependencies, formulaDependency)	
		}
		tx.Where(&FormulaDependencies{
			CategoryID: uint(params.CategoryID),
			TargetAttributeID: uint(params.TargetAttributeID),
		}).Delete(&FormulaDependencies{})

		result := tx.Create(&formulaDependencies)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetFormulas(ctx context.Context, categoryIds []int) []models.FormulasResult {
	var formulas []models.FormulasResult
	s.DB.Model(&Formulas{}).
		Select(`category_id as "categoryId", target_attribute_id as "targetAttributeId", expression`).
		Where("category_id IN ?", categoryIds).
		Scan(&formulas)
	return formulas
}

func (s *Store) GetTopologicalSorting(ctx context.Context, categoryIds []int) ([]models.TopologicalSortResult, error) {
	var topologicalSorting []models.TopologicalSortResult
	err := s.DB.Model(&CategoryAttributeAssignment{}).
		Select(`category_id as "categoryId", attribute_id as "attributeId", topological_sort_order as "topologicalSortOrder"`).
		Where("category_id IN ?", categoryIds).
		Order("category_id, topological_sort_order").
		Scan(&topologicalSorting).Error
	if err != nil {
		return []models.TopologicalSortResult{}, err
	}
	return topologicalSorting, nil
}

func (s *Store) GetProductData(ctx context.Context, productIds []string) ([]models.ProductDatasResult, error) {
	var productDatas []models.ProductDatasResult
	
	err := s.DB.Raw(`
		WITH product_data AS (
			SELECT 
				id, 
				data, 
				attribute_id, 
				category_id 
			FROM 
				products 
			WHERE 
				id IN ?
		)
		SELECT 
			product_data.id, 
			product_data.data, 
			attributes.id AS "attributeId", 
			attributes.name AS "attributeName", 
			attributes.data_type AS "dataType",
			product_data.category_id AS "categoryId"
		FROM 
			product_data 
			JOIN attributes ON product_data.attribute_id = attributes.id
	`, productIds).Scan(&productDatas).Error
	if err != nil{
		return []models.ProductDatasResult{}, err
	}
	return productDatas, nil
}

func (s *Store) GetProductList(ctx context.Context) ([]models.ProductListResult, error) {
	var productList []models.ProductListResult
	err := s.DB.Raw(`
		WITH product_data AS (
			SELECT 
				id, 
				data, 
				attribute_id, 
				category_id 
			FROM 
				products 
			WHERE 
				attribute_id = 0
		)
		SELECT 
			product_data.id, 
			product_data.data as name, 
			categories.path as "categoryPath",
			categories.id as "categoryId"
		FROM 
			product_data 
			JOIN categories ON product_data.category_id = categories.id
	`).Scan(&productList).Error
	if err != nil{
		return []models.ProductListResult{}, err
	}
	return productList, nil
}

func (s *Store) GetFormulasList(ctx context.Context) ([]models.FormulasListResult, error) {
	var formulaList []models.FormulasListResult
	err := s.DB.Raw(`
		select 
			c.id as "categoryId", 
			f.target_attribute_id as "targetAttributeId", 
			c.path as "categoryName", 
			a.name as "targetAttributeName", 
			f.expression as "formula" 
		from 
			formulas f 
			join attributes a on f.target_attribute_id = a.id 
			join categories c on f.category_id = c.id

	`).Scan(&formulaList).Error
	if err != nil{
		return []models.FormulasListResult{}, err
	}
	return formulaList, nil
}


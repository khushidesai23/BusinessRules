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

func NewStore(db *gorm.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) CreateCategory(ctx context.Context, name string) error {
	fmt.Println("Hello2")
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		category := Category{
			Path: name,
		}
		fmt.Println("Hello3")
		err := tx.Save(&category).Error
		fmt.Println("err: ", err)
		if err != nil {
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
				CategoryID:  uint(catID),
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
	var nameValues []string
	var attrValues []string

	// Separate product names (attributeId=0) from attribute values (attributeId>0)
	// This distinction is critical because:
	// - Product names (attributeId=0) are stored in the Name column
	// - Regular attributes are stored in the Data column
	// - Product names don't have corresponding Attribute records (avoiding FK violations)
	for _, data := range datas {
		if data.AttributeID == 0 {
			// Product name entry - will include attributeId=0
			nameValues = append(nameValues, fmt.Sprintf("('%s', %d, '%s')", data.ID, data.CategoryID, data.Data))
		} else {
			// Attribute value entry - includes attributeId
			attrValues = append(attrValues, fmt.Sprintf("('%s', %d, %d, '%s')", data.ID, data.CategoryID, data.AttributeID, data.Data))
		}
	}

	// Upsert product names (id, category_id, name)
	if len(nameValues) > 0 {
		// include attribute_id = 0 for name rows so they satisfy NOT NULL and primary key
		// nameValues currently formatted as ('id', category_id, 'name') — rewrite with attribute_id = 0
		var nameRows []string
		for _, v := range nameValues {
			// v is like ( 'id', 6, 'pen' ) -> insert attribute_id = 0 after category_id
			// we'll transform by injecting , 0 after the second comma position
			// simpler: rebuild from original data by splitting on comma
			parts := strings.SplitN(v, ",", 3)
			if len(parts) == 3 {
				// parts[0]="('id'", parts[1]=" 6", parts[2]=" 'pen')"
				newRow := fmt.Sprintf("%s,%s,0,%s", strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]))
				nameRows = append(nameRows, newRow)
			} else {
				// fallback: append with attribute_id 0 manually
				nameRows = append(nameRows, strings.Replace(v, ")", ", 0)", 1))
			}
		}
		nameQuery := strings.Join(nameRows, ",")
		q := fmt.Sprintf(`
			INSERT INTO products (id, category_id, attribute_id, name)
			VALUES %s
			ON CONFLICT ON CONSTRAINT products_pkey
			DO UPDATE SET name = EXCLUDED.name
		`, nameQuery)
		if err := s.DB.Exec(q).Error; err != nil {
			return err
		}
	}

	// Upsert attribute values (id, category_id, attribute_id, data)
	// These rows have attributeId > 0 and store actual attribute values in the Data column
	if len(attrValues) > 0 {
		queryStr := strings.Join(attrValues, ",")
		query := fmt.Sprintf(`
			INSERT INTO products (id, category_id, attribute_id, data)
			VALUES
				%s
			ON CONFLICT (id, category_id, attribute_id)
			DO UPDATE SET
				data = EXCLUDED.data
		`, queryStr)
		if err := s.DB.Exec(query).Error; err != nil {
			return err
		}
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

func (s *Store) GetAttributesIdDataMap(ctx context.Context, categoryIds []int) (map[int]Attribute, error) {
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
		for index, value := range params.TopologicallySortedAttributeIDs {
			gorm.G[CategoryAttributeAssignment](tx).Where("category_id = ? AND attribute_id = ?", params.CategoryID, value).Update(ctx, "topological_sort_order", index)
		}
		tx.Save(&Formulas{
			CategoryID:        uint(params.CategoryID),
			Expression:        params.Formula,
			TargetAttributeID: uint(params.TargetAttributeID),
		})
		var formulaDependencies []FormulaDependencies
		for _, dependentAttributeId := range params.DependentAttributeIDs {
			formulaDependency := FormulaDependencies{
				CategoryID:           uint(params.CategoryID),
				TargetAttributeID:    uint(params.TargetAttributeID),
				DependentAttributeID: uint(dependentAttributeId),
			}
			formulaDependencies = append(formulaDependencies, formulaDependency)
		}
		tx.Where(&FormulaDependencies{
			CategoryID:        uint(params.CategoryID),
			TargetAttributeID: uint(params.TargetAttributeID),
		}).Delete(&FormulaDependencies{})

		if len(formulaDependencies) > 0 {
			result := tx.Create(&formulaDependencies)
			if result.Error != nil {
				return result.Error
			}
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

	// Enhanced: Query now returns product name as attributeId 0, and other attributes joined to attributes table
	// The query has two parts joined with UNION ALL:
	// 1. Get product names from products.name column (attributeId=0)
	// 2. Get attribute values by joining with attributes table (attributeId>0)
	// This provides a unified interface where product names appear as a special "attribute" with ID 0
	err := s.DB.Raw(`
		SELECT p.id, p.name as data, 0 as "attributeId", 'name' as "attributeName", 'string' as "dataType", p.category_id as "categoryId"
		FROM products p
		WHERE p.id IN ?
		GROUP BY p.id, p.name, p.category_id
		UNION ALL
		SELECT pr.id, pr.data, a.id as "attributeId", a.name as "attributeName", a.data_type as "dataType", pr.category_id as "categoryId"
		FROM products pr
		JOIN attributes a ON pr.attribute_id = a.id
		WHERE pr.id IN ?
	`, productIds, productIds).Scan(&productDatas).Error
	if err != nil {
		return []models.ProductDatasResult{}, err
	}
	return productDatas, nil
}

func (s *Store) GetProductList(ctx context.Context) ([]models.ProductListResult, error) {
	var productList []models.ProductListResult

	// Fixed: Use MAX(p.name) to handle the fact that a product has multiple rows in the products table
	// (one for each attribute: name + all attribute values). We want to get the product name once.
	// GROUP BY is required when selecting MAX() to aggregate correctly by product (id, category)
	// Note: Each product ID appears once per attribute + once for the name row,
	// so we need to group and aggregate to get a unique product record with its name
	err := s.DB.Raw(`
		SELECT p.id, MAX(p.name) as name, c.path as "categoryPath", c.id as "categoryId"
		FROM products p
		JOIN categories c ON p.category_id = c.id
		GROUP BY p.id, c.path, c.id
	`).Scan(&productList).Error
	if err != nil {
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
	if err != nil {
		return []models.FormulasListResult{}, err
	}
	return formulaList, nil
}

// DeleteCategory deletes a category by id (cascades to related records)
func (s *Store) DeleteCategory(ctx context.Context, categoryId int) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", categoryId).Delete(&Category{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteAttribute removes an attribute and related assignments/formulas
func (s *Store) DeleteAttribute(ctx context.Context, attributeId int) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", attributeId).Delete(&Attribute{}).Error; err != nil {
			return err
		}
		// Remove category assignments referencing this attribute
		if err := tx.Where("attribute_id = ?", attributeId).Delete(&CategoryAttributeAssignment{}).Error; err != nil {
			return err
		}
		// Remove formulas that target this attribute
		if err := tx.Where("target_attribute_id = ?", attributeId).Delete(&Formulas{}).Error; err != nil {
			return err
		}
		// Remove formula dependency entries where attribute is dependent
		if err := tx.Where("dependent_attribute_id = ? OR target_attribute_id = ?", attributeId, attributeId).Delete(&FormulaDependencies{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteProduct deletes a product and its attribute rows
func (s *Store) DeleteProduct(ctx context.Context, productId string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", productId).Delete(&Product{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteFormula deletes a saved formula and its dependencies
func (s *Store) DeleteFormula(ctx context.Context, categoryId int, targetAttributeId int) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("category_id = ? AND target_attribute_id = ?", categoryId, targetAttributeId).Delete(&Formulas{}).Error; err != nil {
			return err
		}
		if err := tx.Where("category_id = ? AND target_attribute_id = ?", categoryId, targetAttributeId).Delete(&FormulaDependencies{}).Error; err != nil {
			return err
		}
		return nil
	})
}

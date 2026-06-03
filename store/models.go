package storage

import (
	"log"
	"time"
)

type JSONB map[string]any

// Attribute represents a category attribute (e.g., color, size, price)
type Attribute struct {
	ID        uint       `gorm:"primaryKey; autoIncrement" json:",omitempty"`
	Name      string     `gorm:"uniqueIndex;not null" validate:"required"`
	DataType  string     `gorm:"not null" validate:"required"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:",omitempty"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:",omitempty"`
	DeletedAt *time.Time `gorm:"index" json:",omitempty"`
}

// Formulas stores expression-based calculations for specific attributes within a category
// Fixed: Corrected GORM foreign key references and relationships
// This ensures GORM generates correct SQL FK constraints
type Formulas struct {
	CategoryID        uint      `gorm:"primaryKey"`
	Category          Category  `gorm:"foreignKey:CategoryID; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TargetAttributeID uint      `gorm:"primaryKey"`
	TargetAttribute   Attribute `gorm:"foreignKey:TargetAttributeID; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Expression        string
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:",omitempty"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:",omitempty"`
	DeletedAt         *time.Time `gorm:"index" json:",omitempty"`
}

// FormulaDependencies tracks which attributes a formula depends on
// Simplified: Removed explicit GORM relationship fields (Category, TargetAttribute, DependentAttribute)
type FormulaDependencies struct {
	CategoryID           uint       `gorm:"primaryKey"`
	TargetAttributeID    uint       `gorm:"primaryKey"`
	DependentAttributeID uint       `gorm:"primaryKey"`
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:",omitempty"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:",omitempty"`
	DeletedAt            *time.Time `gorm:"index" json:",omitempty"`
}

// Category represents a hierarchical category grouping (e.g., Electronics > Phones)
// Fixed: Corrected "grom" typo to "gorm" in primary key tag
type Category struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:",omitempty"`
	Path       string      `gorm:"uniqueIndex; not null"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:",omitempty"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime" json:",omitempty"`
	DeletedAt  *time.Time  `gorm:"index" json:",omitempty"`
	Attributes []Attribute `gorm:"many2many:category_attribute_assignments;"`
}

// CategoryAttributeAssignment represents the join table between Category and Attribute
type CategoryAttributeAssignment struct {
	CategoryID           uint      `gorm:"primaryKey"`
	Category             Category  `gorm:"foreignKey:CategoryID; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE"`
	AttributeID          uint      `gorm:"primaryKey"`
	Attribute            Attribute `gorm:"foreignKey:AttributeID; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE"`
	TopologicalSortOrder uint
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:",omitempty"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:",omitempty"`
	DeletedAt            *time.Time `gorm:"index" json:",omitempty"`
}

// Product stores product data with attributes and a dedicated name field
// Fixed & Enhanced:
// - Added Name field: Stores product name separately, reducing data duplication. This allows products to have a human-readable name (attributeId=0)
// - ProductID, CategoryID and AttributeID form composite primary key. Note: AttributeID can be 0 for product name entries (no FK constraint for 0)
// - Fixed foreign key references to properly cascade updates/deletes
type Product struct {
	ID          string   `gorm:"primaryKey" json:",omitempty"`
	CategoryID  uint     `gorm:"primaryKey"`
	Category    Category `gorm:"foreignKey:CategoryID; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AttributeID uint     `gorm:"primaryKey"`
	Name        string   `gorm:"type:text" json:",omitempty"`
	Data        string
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:",omitempty"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:",omitempty"`
	DeletedAt   *time.Time `gorm:"index" json:",omitempty"`
}

type ApiResponse struct {
	Message string `json:"message"`
	Data    []any  `json:"data"`
}

func AutoMigrate() error {
	// Migrate basic tables with simple relationships
	if err := DB.AutoMigrate(&Attribute{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&Category{}); err != nil {
		return err
	}

	// CategoryAttributeAssignment is a manual join table - must be created explicitly
	if !DB.Migrator().HasTable(&CategoryAttributeAssignment{}) {
		if err := DB.Migrator().CreateTable(&CategoryAttributeAssignment{}); err != nil {
			log.Printf("AutoMigrate: failed to create category_attribute_assignments: %v", err)
			return err
		}
	}

	// Ensure columns exist on join table (added after initial schema)
	// TopologicalSortOrder was added later, so check and add if missing
	if DB.Migrator().HasTable(&CategoryAttributeAssignment{}) {
		if !DB.Migrator().HasColumn(&CategoryAttributeAssignment{}, "TopologicalSortOrder") {
			if err := DB.Migrator().AddColumn(&CategoryAttributeAssignment{}, "TopologicalSortOrder"); err != nil {
				log.Printf("AutoMigrate: failed to add column TopologicalSortOrder: %v", err)
				return err
			}
		}
		// ensure timestamp columns exist (use raw SQL to avoid GORM naming mismatches)
		if err := DB.Exec("ALTER TABLE category_attribute_assignments ADD COLUMN IF NOT EXISTS created_at timestamptz").Error; err != nil {
			log.Printf("AutoMigrate: failed to add column created_at: %v", err)
			return err
		}
		if err := DB.Exec("ALTER TABLE category_attribute_assignments ADD COLUMN IF NOT EXISTS updated_at timestamptz").Error; err != nil {
			log.Printf("AutoMigrate: failed to add column updated_at: %v", err)
			return err
		}
		if err := DB.Exec("ALTER TABLE category_attribute_assignments ADD COLUMN IF NOT EXISTS deleted_at timestamptz").Error; err != nil {
			log.Printf("AutoMigrate: failed to add column deleted_at: %v", err)
			return err
		}
	}

	// Migrate formula-related tables
	if err := DB.AutoMigrate(&Formulas{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&FormulaDependencies{}); err != nil {
		return err
	}

	// Migrate Product table (newly enhanced with Name field)
	if err := DB.AutoMigrate(&Product{}); err != nil {
		return err
	}
	// ensure name column exists for product name storage
	if DB.Migrator().HasTable(&Product{}) {
		if !DB.Migrator().HasColumn(&Product{}, "Name") {
			if err := DB.Migrator().AddColumn(&Product{}, "Name"); err != nil {
				log.Printf("AutoMigrate: failed to add column Name to products: %v", err)
				return err
			}
		}
	}
	return nil
}

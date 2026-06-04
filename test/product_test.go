package test

import (
	"bytes"
	"calculationengine/models"
	"calculationengine/router"
	storage "calculationengine/store"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateProduct tests product creation endpoint
func TestCreateProduct(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	router.Api()

	// First, create a category and attribute for the product
	categoryReq := models.CreateCategoryRequest{Name: "Electronics"}
	body, _ := json.Marshal(categoryReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/category/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test category: %s", w.Body.String())
	}

	// Get category ID from database
	var category storage.Category
	db.First(&category, "path = ?", "Electronics")

	// Create an attribute
	attrReq := models.CreateAttributeRequest{Name: "Color", DataType: "string"}
	body, _ = json.Marshal(attrReq)
	req = httptest.NewRequest(http.MethodPost, "/v1/attribute/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test attribute: %s", w.Body.String())
	}

	// Get attribute ID from database
	var attr storage.Attribute
	db.First(&attr, "name = ?", "Color")

	// Assign attribute to category
	assignReq := models.ChangeCategoryAttributeAssignmentRequest{}
	assignReq.Assign.CategoryIDs = []int{int(category.ID)}
	assignReq.Assign.AttributeIDs = []int{int(attr.ID)}
	body, _ = json.Marshal(assignReq)
	req = httptest.NewRequest(http.MethodPost, "/v1/assignment/change", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	// Test case 1: Create product with product name (attributeId=0)
	t.Run("CreateProduct_WithName", func(t *testing.T) {
		productReq := models.CreateProductRequest{
			CategryID: int(category.ID),
			ProductData: []struct {
				AttributeID int    `json:"attributeId" validate:"required"`
				Value       string `json:"value" validate:"required"`
			}{
				{AttributeID: 0, Value: "Awesome Phone"}, // Product name
				{AttributeID: int(attr.ID), Value: "Black"}, // Attribute value
			},
		}

		body, _ := json.Marshal(productReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/product/upsert", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var response storage.ApiResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Message != "success" {
			t.Errorf("Expected success message, got: %s", response.Message)
		}

		LogTest(t, "Product created successfully with name and attributes")
	})

	// Test case 2: Create product without product name (should fail)
	t.Run("CreateProduct_MissingName", func(t *testing.T) {
		productReq := models.CreateProductRequest{
			CategryID: int(category.ID),
			ProductData: []struct {
				AttributeID int    `json:"attributeId" validate:"required"`
				Value       string `json:"value" validate:"required"`
			}{
				{AttributeID: int(attr.ID), Value: "Red"}, // Only attribute, no name
			},
		}

		body, _ := json.Marshal(productReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/product/upsert", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		var response storage.ApiResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		// Should fail because product name (attributeId=0) is required for new products
		if response.Message == "success" {
			t.Errorf("Expected product creation to fail without product name, but succeeded")
		}

		LogTest(t, "Correctly rejected product without name")
	})
}

// TestGetProductList tests retrieving products
func TestGetProductList(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	router.Api()

	// Setup: Create category, attribute, and product
	categoryReq := models.CreateCategoryRequest{Name: "Phones"}
	body, _ := json.Marshal(categoryReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/category/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	var category storage.Category
	db.First(&category, "path = ?", "Phones")

	attrReq := models.CreateAttributeRequest{Name: "Price", DataType: "float"}
	body, _ = json.Marshal(attrReq)
	req = httptest.NewRequest(http.MethodPost, "/v1/attribute/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	var attr storage.Attribute
	db.First(&attr, "name = ?", "Price")

	assignReq := models.ChangeCategoryAttributeAssignmentRequest{}
	assignReq.Assign.CategoryIDs = []int{int(category.ID)}
	assignReq.Assign.AttributeIDs = []int{int(attr.ID)}
	body, _ = json.Marshal(assignReq)
	req = httptest.NewRequest(http.MethodPost, "/v1/assignment/change", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to assign attribute to category: %s", w.Body.String())
	}

	// Create a product
	productReq := models.CreateProductRequest{
		CategryID: int(category.ID),
		ProductData: []struct {
			AttributeID int    `json:"attributeId" validate:"required"`
			Value       string `json:"value" validate:"required"`
		}{
			{AttributeID: 0, Value: "iPhone 15"},
			{AttributeID: int(attr.ID), Value: "999.99"},
		},
	}

	body, _ = json.Marshal(productReq)
	req = httptest.NewRequest(http.MethodPost, "/v1/product/upsert", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test product: %s", w.Body.String())
	}

	// Test get product list
	t.Run("GetProductList_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/product/get-all", nil)
		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		if w.Code != http.StatusOK && w.Code != http.StatusCreated {
			t.Errorf("Expected status 200 or 201, got %d. Response: %s", w.Code, w.Body.String())
		}

		var response struct {
			Message string                       `json:"message"`
			Data    []models.ProductListResult `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Message != "success" {
			t.Errorf("Expected success message, got: %s", response.Message)
		}

		if len(response.Data) == 0 {
			t.Errorf("Expected at least 1 product in list")
			return
		}

		// Verify product name is included
		if response.Data[0].Name != "iPhone 15" {
			t.Errorf("Expected product name 'iPhone 15', got '%s'", response.Data[0].Name)
		}

		LogTest(t, "Retrieved %d products, first product: %s", len(response.Data), response.Data[0].Name)
	})
}

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

// TestCreateCategory tests category creation endpoint
func TestCreateCategory(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	router.Api()

	// Test case 1: Valid category creation
	t.Run("CreateCategory_Success", func(t *testing.T) {
		requestBody := models.CreateCategoryRequest{
			Name: "Electronics",
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/category/create", bytes.NewBuffer(body))
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

		LogTest(t, "Category created successfully: %s", requestBody.Name)
	})

	// Test case 2: Multiple categories
	t.Run("CreateCategory_Multiple", func(t *testing.T) {
		categories := []string{"Phones", "Laptops", "Accessories"}

		for _, catName := range categories {
			requestBody := models.CreateCategoryRequest{Name: catName}
			body, _ := json.Marshal(requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/category/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.Router.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				t.Errorf("Failed to create category %s: %s", catName, w.Body.String())
			}
		}

		LogTest(t, "Created %d categories", len(categories))
	})

	// Test case 3: Missing required field
	t.Run("CreateCategory_MissingName", func(t *testing.T) {
		requestBody := map[string]string{}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/category/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		LogTest(t, "Validation error correctly caught for missing name")
	})
}

// TestGetAllCategories tests retrieving all categories
func TestGetAllCategories(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	router.Api()

	// Create test categories
	categories := []string{"Books", "Movies", "Music"}
	for _, catName := range categories {
		requestBody := models.CreateCategoryRequest{Name: catName}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/category/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create test category: %s", w.Body.String())
		}
	}

	// Test get all categories
	t.Run("GetAllCategories_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/category/get-all", nil)
		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		if w.Code != http.StatusOK && w.Code != http.StatusCreated {
			t.Errorf("Expected status 200 or 201, got %d. Response: %s", w.Code, w.Body.String())
		}

		var response struct {
			Message string                       `json:"message"`
			Data    []models.GetCategoriesResult `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Message != "success" {
			t.Errorf("Expected success message, got: %s", response.Message)
		}

		if len(response.Data) < len(categories) {
			t.Errorf("Expected at least %d categories, got %d", len(categories), len(response.Data))
		}

		LogTest(t, "Retrieved %d categories", len(response.Data))
	})
}

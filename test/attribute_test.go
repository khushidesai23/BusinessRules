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

// TestCreateAttribute tests the attribute creation endpoint
func TestCreateAttribute(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	// Initialize router
	router.Api()

	// Test case 1: Valid attribute creation
	t.Run("CreateAttribute_Success", func(t *testing.T) {
		requestBody := models.CreateAttributeRequest{
			Name:     "Color",
			DataType: "string",
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/attribute/create", bytes.NewBuffer(body))
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

		LogTest(t, "Attribute created successfully: %s", requestBody.Name)
	})

	// Test case 2: Missing required field
	t.Run("CreateAttribute_MissingField", func(t *testing.T) {
		requestBody := map[string]string{
			"dataType": "string",
			// Missing Name field
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/attribute/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		LogTest(t, "Validation error correctly caught for missing field")
	})

	// Test case 3: Different data types
	dataTypes := []string{"integer", "float", "boolean", "string"}
	for _, dataType := range dataTypes {
		t.Run("CreateAttribute_DataType_"+dataType, func(t *testing.T) {
			requestBody := models.CreateAttributeRequest{
				Name:     "Attr_" + dataType,
				DataType: dataType,
			}

			body, _ := json.Marshal(requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/attribute/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.Router.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				t.Errorf("Failed to create attribute with type %s: %s", dataType, w.Body.String())
			}

			LogTest(t, "Created attribute with type: %s", dataType)
		})
	}
}

// TestGetAllAttributes tests retrieving all attributes
func TestGetAllAttributes(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	router.Api()

	// Create a few attributes first
	attrs := []models.CreateAttributeRequest{
		{Name: "Size", DataType: "string"},
		{Name: "Weight", DataType: "float"},
		{Name: "InStock", DataType: "boolean"},
	}

	for _, attr := range attrs {
		body, _ := json.Marshal(attr)
		req := httptest.NewRequest(http.MethodPost, "/v1/attribute/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.Router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create test attribute: %s", w.Body.String())
		}
	}

	// Now test get all attributes
	req := httptest.NewRequest(http.MethodPost, "/v1/attributes/get-all", nil)
	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("Expected status 200 or 201, got %d. Response: %s", w.Code, w.Body.String())
	}

	var response struct {
		Message string                       `json:"message"`
		Data    []models.GetAttributesResult `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Message != "success" {
		t.Errorf("Expected success message, got: %s", response.Message)
	}

	if len(response.Data) < 3 {
		t.Errorf("Expected at least 3 attributes, got %d", len(response.Data))
	}

	LogTest(t, "Retrieved %d attributes", len(response.Data))
}

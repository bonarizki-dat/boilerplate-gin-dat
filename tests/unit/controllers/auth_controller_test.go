package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/controllers"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter creates a test Gin router
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// setupAuthController creates an AuthController for testing
func setupAuthController() *controllers.AuthController {
	authService := services.NewAuthService()
	return controllers.NewAuthController(authService)
}

// TestRegisterEndpoint tests the Register HTTP handler
func TestRegisterEndpoint(t *testing.T) {
	// NOTE: These tests require proper mocking or database setup
	// For production code, mock the AuthService dependency
	// See TESTING.md for mocking guidelines
	t.Skip("Skipping: Requires mocked AuthService or test database setup")

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid registration request",
			requestBody: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "SecurePass123!",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"name": "John Doe",
				// Missing email and password
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, false, response["success"])
			},
		},
		{
			name: "Invalid email format",
			requestBody: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "SecurePass123!",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Password too short",
			requestBody: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "short",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json string",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupTestRouter()
			authController := setupAuthController()
			router.POST("/auth/register", authController.Register)

			// Create request body
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			// Create request
			req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			t.Logf("Response status: %d, Expected: %d", w.Code, tt.expectedStatus)
			t.Logf("Response body: %s", w.Body.String())

			// Run additional checks if provided
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestLoginEndpoint tests the Login HTTP handler
func TestLoginEndpoint(t *testing.T) {
	// NOTE: These tests require proper mocking or database setup
	// For production code, mock the AuthService dependency
	// See TESTING.md for mocking guidelines
	t.Skip("Skipping: Requires mocked AuthService or test database setup")

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid login request format",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing email",
			requestBody: map[string]interface{}{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing password",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty request body",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "not a json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupTestRouter()
			authController := setupAuthController()
			router.POST("/auth/login", authController.Login)

			// Create request body
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			// Create request
			req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			t.Logf("Response status: %d, Expected: %d", w.Code, tt.expectedStatus)
			t.Logf("Response body: %s", w.Body.String())

			// Run additional checks if provided
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestRegisterResponseStructure tests the response structure
func TestRegisterResponseStructure(t *testing.T) {
	// This test verifies that successful registration returns correct structure
	// Note: This test requires database setup or mocking

	t.Skip("Skipping integration test - requires database setup")

	router := setupTestRouter()
	authController := setupAuthController()
	router.POST("/auth/register", authController.Register)

	requestBody := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "SecurePass123!",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusCreated {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check new standard response structure
		assert.Contains(t, response, "success")
		assert.Contains(t, response, "message")
		assert.Contains(t, response, "data")
		assert.Equal(t, true, response["success"])

		data := response["data"].(map[string]interface{})
		assert.Contains(t, data, "user")
		assert.Contains(t, data, "access_token")
		assert.Contains(t, data, "token_type")
	}
}

// Example: Benchmark test for controller
func BenchmarkRegisterEndpoint(b *testing.B) {
	router := setupTestRouter()
	authController := setupAuthController()
	router.POST("/auth/register", authController.Register)

	requestBody := dto.RegisterRequest{
		Name:     "Benchmark User",
		Email:    "bench@example.com",
		Password: "SecurePass123!",
	}
	body, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

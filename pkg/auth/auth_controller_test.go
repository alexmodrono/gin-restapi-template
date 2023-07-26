package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexmodrono/gin-restapi-template/internal/middlewares"
	"github.com/alexmodrono/gin-restapi-template/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Initialize a new gin router for testing
	router := gin.Default()

	// Sets the errors middleware. If this step is skipped, if a route fails the
	// body will just be empty.
	errors_middleware := middlewares.GetErrorsMiddleware(mocks.NewMockLogger(), router)
	errors_middleware.Setup()

	// Initialize mock logger, mock users service, and mock auth service
	logger := &mocks.MockLogger{}
	usersService := &mocks.MockUsersService{}
	authService := &mocks.MockAuthService{}

	// Create the auth controller for testing
	authController := GetAuthController(logger, authService, usersService)
	// Add the route to the router
	router.POST("/login", authController.Login)

	t.Run("ValidCredentials", func(t *testing.T) {
		// Create a request body with valid login credentials
		requestBody := LoginBody{
			Email:    "user@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		// Create a new HTTP request
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request and record the response
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the response status code and body
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Logged in successfully.", response["message"])
		assert.Equal(t, "mock_jwt_token", response["token"])
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		// Create a request body with valid login credentials
		requestBody := LoginBody{
			Email:    "user@example.com",
			Password: "incorrect_password",
		}
		jsonBody, _ := json.Marshal(requestBody)

		// Create a new HTTP request
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request and record the response
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the response status code and body
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "The password provided is incorrect.", response["error"])
	})
}

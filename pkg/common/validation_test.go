/*
Package Name: common
File Name: validation_test.go
Abstract: Tests for the body validation functions.
Author: Alejandro Modroño <alex@sureservice.es>
Created: 07/26/2023
Last Updated: 07/26/2023

# MIT License

# Copyright 2023 Alejandro Modroño Vara

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidation_ValidateBody_Valid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Initialize a new gin router for testing
	router := gin.New()

	// Define a test struct that mimics the request body
	type TestBody struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	// Define a test request body in JSON format
	requestBodyJSON := `{"name": "John Doe", "email": "john.doe@example.com"}`

	// Create a new request with the test JSON body
	req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(requestBodyJSON))
	assert.NoError(t, err)

	// Set the request Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder for capturing the response
	recorder := httptest.NewRecorder()

	// Define the test handler function that calls ValidateBody
	testHandler := func(ctx *gin.Context) {
		var body TestBody
		errors := Validation.ValidateBody(ctx, &body)
		if errors != nil {
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Validation passed successfully"})
	}

	// Register the testHandler as the handler for the test route
	router.POST("/test", testHandler)

	// Perform the test request
	router.ServeHTTP(recorder, req)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse the response body into a map
	var response map[string]string
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert that the response contains the success message
	assert.Equal(t, "Validation passed successfully", response["message"])
}

func TestValidation_ValidateBody_Invalid(t *testing.T) {
	// Initialize a new gin router for testing
	router := gin.Default()

	// Define a test struct that mimics the request body
	type TestBody struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	// Define a test request body in JSON format with missing required fields
	requestBodyJSON := `{"name": ""}`

	// Create a new request with the test JSON body
	req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(requestBodyJSON))
	assert.NoError(t, err)

	// Set the request Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder for capturing the response
	recorder := httptest.NewRecorder()

	// Define the test handler function that calls ValidateBody
	testHandler := func(ctx *gin.Context) {
		var body TestBody
		errors := Validation.ValidateBody(ctx, &body)
		if errors != nil {
			ctx.JSON(http.StatusBadRequest, errors)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Validation passed successfully"})
	}

	// Register the testHandler as the handler for the test route
	router.POST("/test", testHandler)

	// Perform the test request
	router.ServeHTTP(recorder, req)

	// Assert that the response status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Parse the response body into a slice of ValidationErrorMessage
	var responseErrors map[string][]ValidationErrorMessage
	err = json.Unmarshal(recorder.Body.Bytes(), &responseErrors)
	assert.NoError(t, err)

	// Assert that the response errors contain the expected validation errors
	expectedErrors := map[string][]ValidationErrorMessage{
		"errors": {
			{Field: "name", Message: "This field is required."},
			{Field: "email", Message: "This field is required."},
		},
	}
	assert.Equal(t, expectedErrors, responseErrors)
}

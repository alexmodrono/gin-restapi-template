/*
Package Name: lib
File Name: database.go
Abstract: This file contains functions for validating the body of a request.
Author: Alejandro Modroño <alex@sureservice.es>
Created: 07/22/2023
Last Updated: 07/24/2023

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
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ======== NAMESPACES ========

// validationT is used for creating a namespace
type validationT struct{}

// the Validation namespace
var Validation validationT

// ======== TYPES ========

// ErrorMsg represents an error message returned by the API
// when the validation of the body parameters fails.
type ValidationErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ======== PUBLIC METHODS ========

// BodyIsValid binds the body of a gin request to the given struct,
// and aborts the operation with user-friendly error messages if any
// parameters are missing.
func (validationT) ValidateBody(ctx *gin.Context, body interface{}) *gin.H {
	// Check the Content-Type of the request
	if err := ctx.ShouldBind(body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ValidationErrorMessage, len(ve))
			for i, fe := range ve {
				out[i] = ValidationErrorMessage{
					Field:   getJSONFieldName(reflect.TypeOf(body).Elem(), fe.Field()),
					Message: getValidationErrorMessage(fe),
				}
			}
			return &gin.H{"errors": out}
		}
	}

	return nil
}

// ======== PRIVATE METHODS ========

// Helper function to get the form field name
func getFormFieldName(field string) string {
	// Implement your logic to map the field name as needed for form data.
	// This could be based on your specific naming conventions.
	return field
}

// getJSONFieldName returns the json field tag of a field.
func getJSONFieldName(structType reflect.Type, fieldName string) string {
	field, found := structType.FieldByName(fieldName)
	if !found {
		return fieldName
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}

	parts := strings.Split(jsonTag, ",")
	return parts[0]
}

// GetValidationErrorMessage returns a more user-friendly validation error
func getValidationErrorMessage(error validator.FieldError) string {
	switch error.Tag() {
	case "required":
		return "This field is required."
	case "lte":
		return "This field should be less than or equal to " + error.Param() + "."
	case "gte":
		return "This field should be greater than or equal to " + error.Param() + "."
	case "email":
		return "Please enter a valid email address."
	case "eqfield":
		return "Must be equal to " + error.Param() + "."
	}
	return error.Tag()
}

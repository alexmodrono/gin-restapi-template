/*
Package Name: middlewares
File Name: auth_middleware.go
Abstract: The middleware for protecting routes.

Author: Alejandro Modroño <alex@sureservice.es>
Created: 07/08/2023
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
package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alexmodrono/gin-restapi-template/pkg/interfaces"
	"github.com/alexmodrono/gin-restapi-template/pkg/lib"
	"github.com/gin-gonic/gin"
)

// ======== TYPES ========

// AuthMiddleware middleware for authentication
type AuthMiddleware struct {
	service interfaces.AuthService
	logger  lib.Logger
}

// ======== PUBLIC METHODS ========

// GetAuthMiddleware returns the auth middleware
func GetAuthMiddleware(
	logger lib.Logger,
	service interfaces.AuthService,
) AuthMiddleware {
	return AuthMiddleware{
		service: service,
		logger:  logger,
	}
}

// Setup sets up jwt auth middleware
func (middleware AuthMiddleware) Setup() {}

// Handler handles the middleware's functionality
func (middleware AuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieve the Authorization header from the request
		authHeader := ctx.Request.Header.Get("Authorization")
		authHeaderSplit := strings.Split(authHeader, " ")

		if len(authHeaderSplit) != 2 || strings.ToLower(authHeaderSplit[0]) != "bearer" {
			middleware.logger.Info("Tried to access protected route without credentials.")
			// If the Authorization header is missing or does not start with "Bearer",
			// return an HTTP 401 Unauthorized response or handle the error appropriately.
			ctx.AbortWithError(
				http.StatusUnauthorized,
				errors.New("An access token is required for accessing this data."),
			)
			return
		}

		// Extract the token from the Authorization header
		token := authHeaderSplit[1]
		// Check the validity of the token using the authentication service
		id, err := middleware.service.CheckToken(token)
		if err != nil {
			// If there is an error in token verification, return an internal server error
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Set the authenticated user's ID in the context for downstream handlers to access
		ctx.Set("id", id)
		ctx.Next()
		return

	}
}

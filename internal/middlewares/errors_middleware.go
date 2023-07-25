/*
Package Name: middlewares
File Name: errors_middleware.go
Abstract: The error middleware for better error handling.

Author: Alejandro Modroño <alex@sureservice.es>
Created: 07/12/2023
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
	"github.com/alexmodrono/gin-restapi-template/pkg/lib"
	"github.com/gin-gonic/gin"
)

// ======== TYPES ========

// AuthMiddleware middleware for authentication
type ErrorsMiddleware struct {
	logger *lib.Logger
	router *lib.Router
}

// ======== PUBLIC METHODS ========

// GetAuthMiddleware returns the auth middleware
func GetErrorsMiddleware(
	logger *lib.Logger,
	router *lib.Router,
) ErrorsMiddleware {
	return ErrorsMiddleware{
		logger: logger,
		router: router,
	}
}

// Setup sets up errors middleware
func (middleware ErrorsMiddleware) Setup() {
	middleware.logger.Info("Setting up [ERRORS] middleware")
	middleware.router.Use(func(ctx *gin.Context) {
		ctx.Next()

		// if any of the routes abort with an error, it
		// will be catched here and displayed to the user.
		for _, err := range ctx.Errors {
			middleware.logger.Error("An error ocurred:", err)
			ctx.JSON(-1, err)
		}
	})
}

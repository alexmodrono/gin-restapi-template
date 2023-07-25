/*
Package Name: middlewares
File Name: cors_middleware.go
Abstract: The cors middleware for implementing CORS.

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
	"os"

	"github.com/alexmodrono/gin-restapi-template/pkg/lib"
	cors "github.com/rs/cors/wrapper/gin"
)

// ======== TYPES ========

// CorsMiddleware middleware for cors
type CorsMiddleware struct {
	router *lib.Router
	logger *lib.Logger
}

// ======== PUBLIC METHODS ========

// NewCorsMiddleware creates new cors middleware
func GetCorsMiddleware(router *lib.Router, logger *lib.Logger) CorsMiddleware {
	return CorsMiddleware{
		router: router,
		logger: logger,
	}
}

// Setup sets up cors middleware
func (middleware CorsMiddleware) Setup() {
	middleware.logger.Info("Setting up [CORS] middleware")

	debug := os.Getenv("ENVIRONMENT") == "development"
	middleware.router.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
		Debug:            debug,
	}))
}

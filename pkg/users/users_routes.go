/*
Package Name: users
File Name: users_routes.go
Abstract: The file containing all the user routes.

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
package users

import (
	"github.com/alexmodrono/gin-restapi-template/internal/middlewares"
	"github.com/alexmodrono/gin-restapi-template/pkg/lib"
)

// ======== TYPES ========

// UsersRoutes struct
type UsersRoutes struct {
	logger          *lib.Logger
	router          *lib.Router
	usersController UsersController
	authMiddleware  middlewares.AuthMiddleware
}

// ======== PUBLIC METHODS ========

// Returns a UserRoutes struct.
func SetUsersRoutes(
	logger *lib.Logger,
	router *lib.Router,
	usersController UsersController,
	authMiddleware middlewares.AuthMiddleware,
) UsersRoutes {
	return UsersRoutes{
		logger:          logger,
		router:          router,
		usersController: usersController,
		authMiddleware:  authMiddleware,
	}
}

// Setup the user routes
func (route UsersRoutes) Setup() {
	route.logger.Info("Setting up [USERS] routes.")
	api := route.router.Group("/users").Use(route.authMiddleware.Handler())
	{
		api.GET("/", route.usersController.GetAll)
		api.GET("/:id", route.usersController.Get)
	}
}

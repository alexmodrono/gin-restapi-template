/*
Package Name: users
File Name: users_controller.go
Abstract: The user controller for performing operations after a route is called.

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
	"errors"
	"net/http"
	"strconv"

	"github.com/alexmodrono/gin-restapi-template/pkg/lib"
	"github.com/gin-gonic/gin"
)

// ======== TYPES ========

// UsersController data type
type UsersController struct {
	// service domains.UserService
	logger  lib.Logger
	service UsersRepository
}

// ======== METHODS ========

// Creates a new user controller and exposes its routes
// to the router.
func GetUsersController(logger lib.Logger, service UsersRepository) UsersController {
	return UsersController{
		logger:  logger,
		service: service,
	}
}

func (controller UsersController) Get(ctx *gin.Context) {
	// Get the id from the context
	idParam := ctx.Param("id")
	controller.logger.Info("[GET] Getting user with id", idParam)

	// ======== TYPE CONVERSION ========
	// Convert the id from string to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("The id must be an int."))
		return
	}

	// ======== RETRIEVE USER ========
	internalUser, err := controller.service.GetUserById(id)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// The controller.service.GetUser(int) function returns a models.InternalUser
	// struct, which contains the password. To avoid exposing this data to the
	// end user, we must convert the internal user to a public user as follows:
	publicUser := internalUser.ToPublic()

	// We can now return the user
	ctx.JSON(http.StatusOK, publicUser)
}

func (controller UsersController) GetAll(ctx *gin.Context) {
	controller.logger.Info("[GET] Getting all users.")

	// ======== RETRIEVE USER ========
	internalUsers, err := controller.service.GetUsers()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// The controller.service.GetUser(int) function returns a models.InternalUser
	// struct, which contains the password. To avoid exposing this data to the
	// end user, we must convert the internal user to a public user as follows:
	publicUsers := make([]PublicUser, len(internalUsers))
	for user := range internalUsers {
		publicUsers[user] = internalUsers[user].ToPublic()
	}

	// We can now return the user
	ctx.JSON(http.StatusOK, publicUsers)
}

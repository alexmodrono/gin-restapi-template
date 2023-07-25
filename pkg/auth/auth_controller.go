/*
Package Name: auth
File Name: auth_controller.go
Abstract: The controller for everything related with authenticating users.

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
package auth

import (
	"errors"
	"net/http"

	"github.com/alexmodrono/gin-restapi-template/pkg/common"
	"github.com/alexmodrono/gin-restapi-template/pkg/interfaces"
	"github.com/alexmodrono/gin-restapi-template/pkg/lib"
	"github.com/alexmodrono/gin-restapi-template/pkg/users"
	"github.com/gin-gonic/gin"
)

// ======== TYPES ========

// AuthController struct
type AuthController struct {
	logger       *lib.Logger
	service      interfaces.AuthService
	usersService users.UsersService
}

type LoginBody struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type SignupBody struct {
	Username        string `json:"username" form:"username" binding:"required,alpha"`
	Email           string `json:"email" form:"email" binding:"required,email"`
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"`
}

// ======== METHODS ========

// GetAuthController retrieves a new auth controller.
func GetAuthController(
	logger *lib.Logger,
	service interfaces.AuthService,
	usersService users.UsersService,
) AuthController {
	return AuthController{
		logger:       logger,
		service:      service,
		usersService: usersService,
	}
}

// SignIn signs in user
func (controller AuthController) Login(ctx *gin.Context) {
	controller.logger.Info("[POST] Login route.")

	// ======== VALIDATE PARAMETERS ========
	// Initilize an empty DTO that represents the parameters
	// this route expects.
	body := LoginBody{}

	err := ctx.Request.ParseForm()
	if err != nil {
		// Handle the error if parsing fails.
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Validate the body and, if successful, assign the
	// contents to the DTO.
	if errors := common.Validation.ValidateBody(ctx, &body); errors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errors)
		return
	}

	// ======== CHECK CREDENTIALS ========
	// Retrieve the user from the database by the email.
	user, err := controller.usersService.GetUserByEmail(body.Email)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check whether the password is correct using the hasher's
	// compare function.
	matches, err := common.Hasher.Compare(body.Password, user.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if matches {
		// Create a JWT token for the user with the subject.
		token, err := controller.service.CreateToken(user.ID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// And, finally, return the token.
		ctx.JSON(200, gin.H{
			"message": "Logged in successfully.",
			"token":   token,
		})
		return
	}

	ctx.AbortWithError(http.StatusUnauthorized, errors.New("The password provided is incorrect."))
}

// Register registers user
func (controller AuthController) Signup(ctx *gin.Context) {
	controller.logger.Info("[POST] Signup route.")

	// ======== VALIDATE PARAMETERS ========
	// Initilize an empty DTO that represents the parameters
	// this route expects.
	body := SignupBody{}

	// Validate the body and, if successful, assign the
	// contents to the DTO.
	if errors := common.Validation.ValidateBody(ctx, &body); errors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errors)
		return
	}

	// ======== CREATE USER ========

	// Retrieve the user from the database by the email.
	id, err := controller.usersService.CreateUser(body.Email, body.Username, body.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Create a JWT token for the user with the subject.
	token, err := controller.service.CreateToken(*id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// And, finally, return the token.
	ctx.JSON(200, gin.H{
		"message": "User signed-up successfully.",
		"token":   token,
	})
}

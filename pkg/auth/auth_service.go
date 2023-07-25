/*
Package Name: auth
File Name: auth_service.go
Abstract: The service for verifying and creating JWTs.

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
	"os"
	"time"

	"github.com/alexmodrono/gin-restapi-template/pkg/interfaces"
	"github.com/alexmodrono/gin-restapi-template/pkg/lib"
	"github.com/golang-jwt/jwt/v5"
)

// ======== TYPES ========
// AuthService service layer
type AuthService struct {
	db *lib.Database
}

// ======== METHODS ========
// GetUserService returns the user service.
func GetAuthService(logger *lib.Logger, db *lib.Database) interfaces.AuthService {
	return AuthService{
		db: db,
	}
}

// CheckToken checks whether the token is correct and returns the subject, which
// in the case of our API is supposed to be the id of the user.
func (service AuthService) CheckToken(tokenString string) (*int32, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		subFloat := claims["sub"].(float64)
		sub := int32(subFloat)
		return &sub, nil
	}

	return nil, err
}

// CreateToken creates jwt auth token
func (service AuthService) CreateToken(id int32) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"iat": time.Now().Unix(),
		"exp": time.Now().AddDate(0, 0, 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

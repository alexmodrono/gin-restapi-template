/*
Package Name: users
File Name: users_controller.go
Abstract: A representation of a user in the database.

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
	"time"
)

// ======== TYPES ========

// InternalUser is a struct that represents a user, and it contains its password.
// As its own name suggests, this type should only be used internally.
type InternalUser struct {
	ID        int32
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// PublicUser is basically a user that will be returned by the api. As its own
// name says, it should be used for returning user data publicly.
type PublicUser struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ======== PUBLIC METHODS ========

// Converts an internal user to a public user.
func (self InternalUser) ToPublic() PublicUser {
	return PublicUser{
		ID:        self.ID,
		Username:  self.Username,
		Email:     self.Email,
		CreatedAt: self.CreatedAt,
	}
}

// Creates a new instance of an internal user from data.
func InternalUserFromData(values []interface{}) InternalUser {
	return InternalUser{
		ID:        values[0].(int32),
		Username:  values[1].(string),
		Email:     values[2].(string),
		Password:  values[3].(string),
		CreatedAt: values[4].(time.Time),
	}
}

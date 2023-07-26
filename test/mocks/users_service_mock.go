/*
Package Name: mocks
File Name: users_service_mock.go
Abstract: Interface for mocking the users service in tests.
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
package mocks

import (
	"errors"
	"time"

	"github.com/alexmodrono/gin-restapi-template/pkg/common"
	"github.com/alexmodrono/gin-restapi-template/pkg/users"
)

// Mock UsersService for testing purposes
type MockUsersService struct{}

func (s *MockUsersService) GetUserById(id int) (*users.InternalUser, error) {
	// Mock the GetUserByEmail method to return a test user with a known password
	// for testing the login functionality.
	if id == 1 {
		password, _ := common.Hasher.Hash("password123")
		return &users.InternalUser{
			ID:       1,
			Username: "user",
			Email:    "user@example.com",
			Password: password,
		}, nil
	}
	return nil, errors.New("user not found")
}

func (s *MockUsersService) GetUserByEmail(email string) (*users.InternalUser, error) {
	// Mock the GetUserByEmail method to return a test user with a known password
	// for testing the login functionality.
	if email == "user@example.com" {
		password, _ := common.Hasher.Hash("password123")
		return &users.InternalUser{
			ID:       1,
			Username: "user",
			Email:    "user@example.com",
			Password: password,
		}, nil
	}
	return nil, errors.New("user not found")
}

func (s *MockUsersService) GetUsers() ([]users.InternalUser, error) {
	users := []users.InternalUser{
		{
			ID:        1,
			Username:  "user",
			Email:     "user@example.com",
			Password:  "$argon2id$v=18$m=65536,t=3,p=2$Zm9v$MTIzNDU2",
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			Username:  "user2",
			Email:     "user2@example.com",
			Password:  "$argon2id$v=18$m=65536,t=3,p=2$Zm9v$MTIzNDU2",
			CreatedAt: time.Now(),
		},
	}
	return users, nil
}

func (s *MockUsersService) CreateUser(email, username, password string) (*int32, error) {
	// Mock the CreateUser method to return a test user ID for the signup functionality.
	// You can replace this with any logic to generate a mock user ID for testing.
	userID := int32(1)
	return &userID, nil
}

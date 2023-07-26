/*
Package Name: users
File Name: users_service.go
Abstract: The user service for performing operations upon users in the database.

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
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alexmodrono/gin-restapi-template/pkg/common"
	"github.com/alexmodrono/gin-restapi-template/pkg/lib"
	"github.com/jackc/pgx/v5/pgconn"
)

// ======== TYPES ========

// UsersService service layer
type UsersService struct {
	logger lib.Logger
	db     *lib.Database
}

// ======== PUBLIC METHODS ========

// GetUsersService returns the user service.
func GetUsersService(logger lib.Logger, db *lib.Database) UsersRepository {
	return UsersService{
		logger: logger,
		db:     db,
	}
}

// GetUserById returns a single user with the specified id.
//
// NOTE: This query returns the user with its hashed password, so make sure to convert its value to
// a models.PublicUser struct which omits the password.
func (service UsersService) GetUserById(id int) (*InternalUser, error) {
	service.logger.Info("Retrieving user with id", id)
	return service.getUserByQuery("id", id)
}

// GetUserByEmail returns a single user with the specified email.
//
// NOTE: This query returns the user with its hashed password, so make sure to convert its value to
// a models.PublicUser struct which omits the password.
func (service UsersService) GetUserByEmail(email string) (*InternalUser, error) {
	service.logger.Info("Retrieving user with email", email)
	return service.getUserByQuery("email", email)
}

// GetUsers returns all the users
func (service UsersService) GetUsers() (users []InternalUser, err error) {
	rows, err := service.db.Query(context.Background(), "SELECT * FROM auth.user;")
	service.logger.Info("Retrieving all users.")
	if err != nil {
		service.logger.Fatal("Error while executing query. Err:", err)
		return nil, err
	}

	var results []InternalUser

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			service.logger.Fatal("Error while iterating dataset. Err:", err)
			return nil, err
		}

		// Once the values have been obtained, they need to be converted into Go
		// types.
		results = append(results, InternalUserFromData(values))
	}

	return results, nil
}

// CreateUser inserts a new user in the database
func (service UsersService) CreateUser(email string, username string, password string) (*int32, error) {

	// ======== HASHING THE PASSWORD ========
	hashedPassword, err := common.Hasher.Hash(password)
	if err != nil {
		service.logger.Fatal("An error ocurred while hashing the password:", err)
		return nil, err
	}

	// ======== QUERIES ========
	var id int32
	err = service.db.QueryRow(
		context.Background(),
		`INSERT INTO auth.user VALUES (DEFAULT, $1, $2, $3, $4) RETURNING id;`,
		username,
		email,
		hashedPassword,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return handleError(err, username, email)
	}

	// Return the first user in the result set.
	return &id, nil
}

// TODO: add update and delete operations for users.

// ======== PRIVATE METHODS ========

// Converts an error to a more user-friendly error.
func handleError(err error, username string, email string) (*int32, error) {
	// Check if the error is a PostgreSQL error (*pgconn.PgError)
	// and handle unique constraint violations based on the constraint name.
	if pgerr, ok := err.(*pgconn.PgError); ok {
		if pgerr.ConstraintName == "user_username_unique" {
			// The username already exists, return a specific error message.
			return nil, fmt.Errorf("Username %s is already taken.", username)
		} else if pgerr.ConstraintName == "user_email_unique" {
			// The email already exists, return a specific error message.
			return nil, fmt.Errorf("User with email %s already exists.", email)
		} else {
			// Handle other PostgreSQL errors.
			return nil, fmt.Errorf("Unexpected error while performing operation on user %s: %v\n", email, pgerr)
		}
	}
	// Handle other types of errors (non-PostgreSQL errors).
	return nil, fmt.Errorf("Unexpected error while performing operation on user %s: %v\n", email, err)
}

// getUserByQuery returns a user from the database based on a specific query.
//
// The auth.get_user_by_id($1) SQL function is a custom function designed to retrieve a user with a specific ID.
// Alternatively, you can directly execute a query like:
//
// SELECT u.id, u.username, u.email, u.password, u.created_at FROM auth.user u WHERE u.id = $1;
//
// Similarly, the auth.get_user_by_email($1) function allows you to retrieve a user based on their email.
// You can use the following query as an alternative:
//
// SELECT u.id, u.username, u.email, u.password, u.created_at FROM auth.user u WHERE u.email = $1;
//
// Likewise, the auth.get_user_by_username($1) function retrieves a user based on their username.
// The equivalent query can be used as an alternative:
//
// SELECT u.id, u.username, u.email, u.password, u.created_at FROM auth.user u WHERE u.username = $1;
func (service UsersService) getUserByQuery(queryType string, args ...interface{}) (*InternalUser, error) {
	var query string
	switch queryType {
	case "id":
		query = "SELECT * FROM auth.get_user_by_id($1)"
	case "email":
		query = "SELECT * FROM auth.get_user_by_email($1)"
	case "username":
		query = "SELECT * FROM auth.get_user_by_username($1)"
	default:
		return nil, errors.New("Invalid query type")
	}

	rows, err := service.db.Query(context.Background(), query, args...)
	if err != nil {
		service.logger.Fatal("Error while executing query. Err:", err)
		return nil, err
	}

	// iterate through the rows
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			service.logger.Fatal("Error while iterating dataset. Err:", err)
			return nil, err
		}

		// Once the values have been obtained, they need to be converted into Go types.
		user := InternalUserFromData(values)

		return &user, nil
	}

	var message string
	if val, ok := args[0].(int); ok {
		message = fmt.Sprintf(
			"The user with the %s '%d' could not be found.",
			queryType,
			val,
		)
	} else if val, ok := args[0].(string); ok {
		message = fmt.Sprintf(
			"The user with the %s '%s' could not be found.",
			queryType,
			val,
		)
	} else {
		// Handle the case when args[0] is neither int nor string.
		message = "Invalid value for the user query."
	}

	return nil, errors.New(message)
}

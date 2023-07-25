/*
Package Name: lib
File Name: database.go
Abstract: This file in the 'lib' package contains a method named 'GetDatabase'
that enables asynchronous connection to a PostgreSQL database. Using the 'pgxpool'
library, it establishes a database pool by constructing the connection URL from
environment variables. The method returns the database pool for seamless database interaction
and logs the connection status using a provided logger. It plays a crucial role in
facilitating database connectivity in the 'lib' package.
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
package lib

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ======== TYPES ========

// A type alias for the connection pool
type Database = pgxpool.Pool

// ======== METHODS ========

// GetDatabase returns a database pool to connect to the database asynchronously
func GetDatabase(logger *Logger) *Database {

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	// Create a connection pool to the database using pgxpool
	dbPool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		logger.Fatal("Unable to connect to database: ", err)
		os.Exit(1)
	}

	logger.Info("Connected to the database successfully.")
	// Closes the pool once the function goes out of scope.
	// defer dbPool.Close()

	return dbPool
}

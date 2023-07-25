/*
Package Name: main
File Name: main.go
Abstract: The entry point of the project and the source code of the
executable that will be used for initializing the API.

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
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alexmodrono/gin-restapi-template/internal/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// ======== PRIVATE METHODS ========

// isValidEnvironment checks whether the environment flag is
// a valid environment name.
func isValidEnvironment(environment *string) bool {
	switch *environment {
	case
		"development",
		"production",
		"test":
		return true
	}
	return false
}

// ======== ENTRY POINT ========
func main() {

	//	======== CHECK ENVIRONMENT ========
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	//	The only available environments are "production", "development",
	//	and "test". If any other environment is provided the api should
	//	exit.
	if !isValidEnvironment(environment) {
		fmt.Println("The environment is not valid!")
		os.Exit(1)
	}

	// Set the 'ENVIRONMENT' value to the flag passed so that
	// we can check the state wherever we want.
	os.Setenv("ENVIRONMENT", *environment)

	// ======== CONFIG FILES ========
	// Load the corresponding environment file
	godotenv.Load("configs/.env." + *environment)

	// Load the main file
	godotenv.Load("configs/.env")

	// ======== DISCLAIMER ========
	fmt.Printf("Welcome to %s %s; Written by %s\n", os.Getenv("APP_NAME"), os.Getenv("APP_VERSION"), os.Getenv("APP_AUTHOR"))

	// ======== DEPENDENCY INJECTION ========
	// The api is divided following the next structure:
	// - Bootstrap: bundles all the dependency injection under one `fx.Options` variable for cleaner code.

	// If the environment is production, set the gin
	// environment to 'release'.
	if *environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// fx.NopLogger disables the logger
	fx.New(
		bootstrap.Module,
		fx.NopLogger,
	).Run()
}

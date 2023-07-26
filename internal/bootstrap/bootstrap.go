/*
Package Name: bootstrap
File Name: bootstrap.go
Abstract: Wrapper for invoking all the module's dependencies and starting
the API by loading the essential initial components that allow it to run
and perform more complex tasks.

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
package bootstrap

import (
	"context"
	"fmt"
	"os"

	"github.com/alexmodrono/gin-restapi-template/internal/middlewares"
	"github.com/alexmodrono/gin-restapi-template/pkg/auth"
	"github.com/alexmodrono/gin-restapi-template/pkg/lib"
	"github.com/alexmodrono/gin-restapi-template/pkg/users"
	"go.uber.org/fx"
)

// ======== PRIVATE METHODS ========

// registerHooks registers a lifecycle hook that starts the API and logs a message
// when the app is stopped.
func registerHooks(
	lifecycle fx.Lifecycle,
	router *lib.Router,
	logger lib.Logger,
	routes Routes,
	middlewares middlewares.Middlewares,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				// Log the start of the application with the configured host and port
				logger.Info(
					fmt.Sprintf(
						"Starting application in %s:%s",
						os.Getenv("APP_HOST"),
						os.Getenv("APP_PORT"),
					),
				)

				// ======== SET UP COMPONENTS ========
				// Perform any necessary setup or initialization tasks for the middlewares
				middlewares.Setup()

				// Perform any necessary setup or initialization tasks for the routes
				routes.Setup()

				// Start the router by running it in a separate goroutine
				go router.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))

				return nil
			},
			OnStop: func(ctx context.Context) error {
				// Log the stop of the application and any associated error
				logger.Fatal(
					fmt.Sprintf(
						"Stopping application. Error: %s", ctx.Err(),
					),
				)

				return nil
			},
		},
	)
}

// ======== EXPORTS ========

// Module exports for fx
var Module = fx.Options(
	// Module exports
	lib.Module,
	middlewares.Module,

	// Context exports
	users.Context,
	auth.Context,

	// Bootstrap exports
	fx.Provide(GetRoutes),

	// Methods
	fx.Invoke(registerHooks),
)

// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/go-swagger/go-swagger/examples/composed-auth/restapi/operations"

	models "github.com/go-swagger/go-swagger/examples/composed-auth/models"

	auth "github.com/go-swagger/go-swagger/examples/composed-auth/auth"
	logging "github.com/op/go-logging"
)

//go:generate swagger generate server --target .. --name multiAuthExample --spec ../swagger.yml --principal models.Principal

func configureFlags(api *operations.MultiAuthExampleAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MultiAuthExampleAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})

	// want some color
	format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} [%{module}]%{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	backendFmt := logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format)
	logging.SetLevel(logging.DEBUG, "api")
	logging.SetBackend(backendFmt)
	logger := logging.MustGetLogger("api")

	api.Logger = logger.Infof

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Authorization: Basic" header is set with the Basic scheme
	api.IsRegisteredAuth = func(user string, pass string) (*models.Principal, error) {
		// The header: Authorization: Basic {base64 string} has already been decoded by the runtime as a username:password pair
		api.Logger("IsRegisteredAuth handler called")
		return auth.IsRegistered(user, pass)
	}

	// Applies when the "Authorization: Bearer" header or the "access_token" query is set
	api.HasRoleAuth = func(token string, scopes []string) (*models.Principal, error) {
		// The header: Authorization: Bearer {base64 string} (or ?access_token={base 64 string} param) has already
		// been decoded by the runtime as a token
		api.Logger("HasRoleAuth handler called")
		return auth.HasRole(token, scopes)
	}

	// Applies when the "CustomKeyAsQuery" query is set
	api.IsResellerQueryAuth = func(token string) (*models.Principal, error) {
		api.Logger("ResellerQueryAuth handler called")
		return auth.IsReseller(token)
	}

	// Applies when the "X-Custom-Key" header is set
	api.IsResellerAuth = func(token string) (*models.Principal, error) {
		api.Logger("IsResellerAuth handler called")
		return auth.IsReseller(token)
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	api.AddOrderHandler = operations.AddOrderHandlerFunc(func(params operations.AddOrderParams, principal *models.Principal) middleware.Responder {
		logger.Warningf("AddOrder called with params: %s, and principal: %s", spew.Sdump(params.Order), spew.Sdump(principal))
		return middleware.NotImplemented("operation .AddOrder has not yet been implemented")
	})
	api.GetItemsHandler = operations.GetItemsHandlerFunc(func(params operations.GetItemsParams) middleware.Responder {
		logger.Warningf("GetItems called with NO params and NO principal")
		return middleware.NotImplemented("operation .GetItems has not yet been implemented")
	})
	api.GetOrderHandler = operations.GetOrderHandlerFunc(func(params operations.GetOrderParams, principal *models.Principal) middleware.Responder {
		logger.Warningf("GetOrder called with params: %s, and principal: %s", spew.Sdump(params.OrderID), spew.Sdump(principal))
		return middleware.NotImplemented("operation .GetOrder has not yet been implemented")
	})
	api.GetOrdersForItemHandler = operations.GetOrdersForItemHandlerFunc(func(params operations.GetOrdersForItemParams, principal *models.Principal) middleware.Responder {
		logger.Warningf("GetOrdersForItem called with params: %v, and principal: %v", spew.Sdump(params.ItemID), spew.Sdump(principal))
		return middleware.NotImplemented("operation .GetOrdersForItem has not yet been implemented")
	})
	api.GetAccountHandler = operations.GetAccountHandlerFunc(func(params operations.GetAccountParams, principal *models.Principal) middleware.Responder {
		logger.Warningf("GetAccount called with NO params, and principal: %s", spew.Sdump(principal))
		return middleware.NotImplemented("operation .GetAccount has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

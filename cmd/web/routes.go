package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes(cfg config) http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
	dynamicRequireAuthMiddleware := dynamicMiddleware.Append(app.requireAuthentication)

	mux := pat.New()
	// Home
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	// About
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))

	// Snippet
	mux.Get("/snippet/create", dynamicRequireAuthMiddleware.ThenFunc(app.createSnippedForm))
	mux.Post("/snippet/create", dynamicRequireAuthMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// User
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicRequireAuthMiddleware.ThenFunc(app.logoutUser))
	mux.Get("/user/profile", dynamicRequireAuthMiddleware.ThenFunc(app.profile))
	mux.Get("/user/change-password", dynamicRequireAuthMiddleware.ThenFunc(app.changePasswordForm))
	mux.Post("/user/change-password", dynamicRequireAuthMiddleware.ThenFunc(app.changePassword))

	// Healthcheck
	mux.Get("/ping", http.HandlerFunc(ping))

	// Static
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}

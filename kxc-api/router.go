package api

import (
	"github.com/didil/kubexcloud/kxc-api/handlers"
	mid "github.com/didil/kubexcloud/kxc-api/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// BuildRouter builds the router
func BuildRouter(root *handlers.Root) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(mid.Cors)

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Heartbeat("/ping"))

	authentication := mid.Authentication(root)

	// Routes

	mux.Route("/v1", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			// POST /v1/users/login
			r.Post("/login", root.HandleLoginUser)

			// POST /v1/users
			r.Post("/", root.HandleCreateUser)
		})

		r.With(authentication).Route("/projects", func(r chi.Router) {
			// Get /v1/projects
			r.Get("/", root.HandleListProjects)
			// POST /v1/projects
			r.Post("/", root.HandleCreateProject)

			r.Route("/{project}/apps", func(r chi.Router) {
				// POST /v1/projects/:project/apps/:app/restart
				r.Post("/{app}/restart", root.HandleRestartApp)
				// POST /v1/projects/:project/apps
				r.Post("/", root.HandleCreateApp)
				// GET /v1/projects/:project/apps
				r.Get("/", root.HandleListApps)
				// PUT /v1/projects/:project/apps/:app
				r.Put("/{app}", root.HandleUpdateApp)
			})
		})
	})

	return mux
}

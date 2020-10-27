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
			mux.Post("/login", root.HandleLoginUser)

			// POST /v1/users/create
			mux.Post("/create", root.HandleCreateUser)
		})

		r.With(authentication).Route("/projects", func(r chi.Router) {
			// POST /v1/projects
			r.Post("/", root.HandleCreateProject)

			// POST /v1projects/:project/apps
			r.Post("/{project}/apps", root.HandleCreateApp)
			// GET /v1projects/:project/apps
			r.Get("/{project}/apps", root.HandleListApps)
			// PUT /v1projects/:project/apps/:app
			r.Put("/{project}/apps/{app}", root.HandleUpdateApp)
		})
	})

	return mux
}

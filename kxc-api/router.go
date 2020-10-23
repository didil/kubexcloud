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

	// Routes

	mux.Route("/v1", func(r chi.Router) {
		r.Route("/projects", func(r chi.Router) {
			// POST /projects
			r.Post("/", root.HandleCreateProject)

			// POST /projects/:project/apps
			r.Post("/{project}/apps", root.HandleCreateApp)
		})
	})

	return mux
}

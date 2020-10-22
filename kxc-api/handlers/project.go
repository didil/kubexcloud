package handlers

import (
	"net/http"

	"github.com/didil/kubexcloud/kxc-api/requests"
)

// HandleCreateProject creates a project
func (app *App) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	reqData := &requests.CreateProject{}

	err := readJSON(r, reqData)
	if err != nil {
		app.HandleError(w, r, err)
		return
	}

	err = app.ProjectSvc.Create(r.Context(), reqData)
	if err != nil {
		app.HandleError(w, r, err)
		return
	}

	JSONOk(w, &struct{}{})
}

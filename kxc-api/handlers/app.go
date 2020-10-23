package handlers

import (
	"fmt"
	"net/http"

	"github.com/didil/kubexcloud/kxc-api/requests"
	"github.com/go-chi/chi"
)

// HandleCreateApp creates an app
func (root *Root) HandleCreateApp(w http.ResponseWriter, r *http.Request) {
	projectName := chi.URLParam(r, "project")

	reqData := &requests.CreateApp{}

	err := readJSON(r, reqData)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	// check if the project exists
	project, err := root.ProjectSvc.Get(r.Context(), projectName)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}
	if project == nil {
		root.HandleError(w, r, fmt.Errorf("project not found: %s", projectName))
		return
	}

	err = root.AppSvc.Create(r.Context(), projectName, reqData)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	JSONOk(w, &struct{}{})
}

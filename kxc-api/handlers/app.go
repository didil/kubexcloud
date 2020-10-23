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

// HandleUpdateApp updates an app
func (root *Root) HandleUpdateApp(w http.ResponseWriter, r *http.Request) {
	projectName := chi.URLParam(r, "project")
	appName := chi.URLParam(r, "app")

	reqData := &requests.UpdateApp{}

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

	err = root.AppSvc.Update(r.Context(), projectName, appName, reqData)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	JSONOk(w, &struct{}{})
}

// HandleListApps lists apps
func (root *Root) HandleListApps(w http.ResponseWriter, r *http.Request) {
	projectName := chi.URLParam(r, "project")

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

	respData, err := root.AppSvc.List(r.Context(), projectName)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	JSONOk(w, respData)

}

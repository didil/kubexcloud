package handlers

import (
	"net/http"

	"github.com/didil/kubexcloud/kxc-api/requests"
)

// HandleCreateProject creates a project
func (root *Root) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	userName := r.Context().Value(CtxKey("userName")).(string)

	reqData := &requests.CreateProject{}

	err := readJSON(r, reqData)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	err = root.ProjectSvc.Create(r.Context(), userName, reqData)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	JSONOk(w, &struct{}{})
}

// HandleListProjects lists projects
func (root *Root) HandleListProjects(w http.ResponseWriter, r *http.Request) {
	userName := r.Context().Value(CtxKey("userName")).(string)

	respData, err := root.ProjectSvc.List(r.Context(), userName)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	JSONOk(w, respData)
}

package handlers

import (
	"net/http"

	"github.com/didil/kubexcloud/kxc-api/requests"
)

// HandleCreateProject creates a project
func (root *Root) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	reqData := &requests.CreateProject{}

	err := readJSON(r, reqData)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	err = root.ProjectSvc.Create(r.Context(), reqData)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	JSONOk(w, &struct{}{})
}

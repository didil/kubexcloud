package handlers

import (
	"net/http"
	"os"

	"github.com/didil/kubexcloud/kxc-api/requests"
	"github.com/didil/kubexcloud/kxc-api/responses"
)

// HandleLoginUser login user
func (root *Root) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	reqData := &requests.LoginUser{}
	err := readJSON(r, reqData)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	token, err := root.UserSvc.Login(r.Context(), reqData.Name, reqData.Password)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	respData := &responses.LoginUser{
		Token: token,
	}

	JSONOk(w, respData)
}

// HandleCreateUser creates user
func (root *Root) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	adminKey := r.Header.Get("ADMIN_KEY")

	// authorize to create user
	if adminKey != os.Getenv("ADMIN_KEY") {
		JSONError(w, "action not allowed", http.StatusUnauthorized)
		return
	}

	reqData := &requests.CreateUser{}

	err := readJSON(r, reqData)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	err = root.UserSvc.Create(r.Context(), reqData.Name, reqData.Password)
	if err != nil {
		root.HandleError(w, r, err)
		return
	}

	JSONOk(w, &struct{}{})
}
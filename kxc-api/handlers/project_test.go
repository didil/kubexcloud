package handlers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/didil/kubexcloud/kxc-api"
	"github.com/didil/kubexcloud/kxc-api/handlers"
	"github.com/didil/kubexcloud/kxc-api/requests"
	"github.com/didil/kubexcloud/kxc-api/responses"
	"github.com/didil/kubexcloud/kxc-api/testsupport"
	"github.com/didil/kubexcloud/kxc-api/testsupport/auth"
	"github.com/didil/kubexcloud/kxc-api/testsupport/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProjectTestSuite struct {
	suite.Suite
}

func (suite *ProjectTestSuite) SetupSuite() {
	testsupport.BootstrapTests("../.env.test")
}
func TestProjectTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}

func (suite *ProjectTestSuite) Test_HandleCreateProject_Ok() {
	userName := "test-user"
	token, err := auth.Login(userName)
	suite.NoError(err)

	projectSvc := new(mocks.ProjectSvc)
	root := &handlers.Root{ProjectSvc: projectSvc}

	reqData := &requests.CreateProject{
		Name: "project-a",
	}

	projectSvc.On("Create", mock.AnythingOfType("*context.valueCtx"), userName, reqData).Return(nil)

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, s.URL+"/v1/projects", &b)
	suite.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)

	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))

	respData, err := ioutil.ReadAll(resp.Body)
	suite.NoError(err)
	suite.Equal("{}", string(respData))

	projectSvc.AssertExpectations(suite.T())
}

func (suite *ProjectTestSuite) Test_HandleListProjects_Ok() {
	userName := "test-user"
	token, err := auth.Login(userName)
	suite.NoError(err)

	projectSvc := new(mocks.ProjectSvc)
	root := &handlers.Root{ProjectSvc: projectSvc}

	rawRespData := &responses.ListProject{
		Projects: []responses.ListProjectEntry{
			responses.ListProjectEntry{
				Name: "project-a",
			},
		},
	}

	projectSvc.On("List", mock.AnythingOfType("*context.valueCtx"), userName).Return(rawRespData, nil)

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	req, err := http.NewRequest(http.MethodGet, s.URL+"/v1/projects", nil)
	suite.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)

	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))

	var respData *responses.ListProject
	err = json.NewDecoder(resp.Body).Decode(&respData)
	suite.NoError(err)

	suite.Len(respData.Projects, 1)
	project_1 := respData.Projects[0]
	suite.Equal("project-a", project_1.Name)

	projectSvc.AssertExpectations(suite.T())
}

func (suite *ProjectTestSuite) Test_HandleListProjects_NoAuth() {
	projectSvc := new(mocks.ProjectSvc)
	root := &handlers.Root{ProjectSvc: projectSvc}

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	req, err := http.NewRequest(http.MethodGet, s.URL+"/v1/projects", nil)
	suite.NoError(err)

	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)

	defer resp.Body.Close()
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)

	projectSvc.AssertExpectations(suite.T())
}

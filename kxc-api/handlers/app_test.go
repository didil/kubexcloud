package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type AppTestSuite struct {
	suite.Suite
}

func (suite *AppTestSuite) SetupSuite() {
	testsupport.BootstrapTests("../.env.test")
}
func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func (suite *AppTestSuite) Test_HandleCreateApp_Ok() {
	userName := "test-user"
	token, err := auth.Login(userName)
	suite.NoError(err)

	appSvc := new(mocks.AppSvc)
	projectSvc := new(mocks.ProjectSvc)
	root := &handlers.Root{AppSvc: appSvc, ProjectSvc: projectSvc}

	reqData := &requests.CreateApp{
		Name: "app-a",
	}

	projName := "project-a"
	proj := &responses.Project{
		Name: projName,
	}

	projectSvc.On("Get", mock.AnythingOfType("*context.valueCtx"), userName, projName).Return(proj, nil)
	appSvc.On("Create", mock.AnythingOfType("*context.valueCtx"), projName, reqData).Return(nil)

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, s.URL+"/v1/projects/project-a/apps", &b)
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

	appSvc.AssertExpectations(suite.T())
}

func (suite *AppTestSuite) Test_HandleUpdateApp_Ok() {
	userName := "test-user"
	token, err := auth.Login(userName)
	suite.NoError(err)

	appSvc := new(mocks.AppSvc)
	projectSvc := new(mocks.ProjectSvc)
	root := &handlers.Root{AppSvc: appSvc, ProjectSvc: projectSvc}

	appName := "app-a"
	reqData := &requests.UpdateApp{
		Replicas: 6,
	}

	projName := "project-a"
	proj := &responses.Project{
		Name: projName,
	}

	projectSvc.On("Get", mock.AnythingOfType("*context.valueCtx"), userName, projName).Return(proj, nil)
	appSvc.On("Update", mock.AnythingOfType("*context.valueCtx"), projName, appName, reqData).Return(nil)

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPut, s.URL+fmt.Sprintf("/v1/projects/%s/apps/%s", projName, appName), &b)
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

	appSvc.AssertExpectations(suite.T())
}

func (suite *AppTestSuite) Test_HandleListApps_Ok() {
	userName := "test-user"
	token, err := auth.Login(userName)
	suite.NoError(err)

	appSvc := new(mocks.AppSvc)
	projectSvc := new(mocks.ProjectSvc)
	root := &handlers.Root{AppSvc: appSvc, ProjectSvc: projectSvc}

	rawRespData := &responses.ListApp{
		Apps: []responses.ListAppEntry{
			responses.ListAppEntry{
				Name: "app-a",
			},
		},
	}

	projName := "project-a"
	proj := &responses.Project{
		Name: projName,
	}

	projectSvc.On("Get", mock.AnythingOfType("*context.valueCtx"), userName, projName).Return(proj, nil)
	appSvc.On("List", mock.AnythingOfType("*context.valueCtx"), projName).Return(rawRespData, nil)

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	req, err := http.NewRequest(http.MethodGet, s.URL+fmt.Sprintf("/v1/projects/%s/apps", projName), nil)
	suite.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)

	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))

	var respData *responses.ListApp
	err = json.NewDecoder(resp.Body).Decode(&respData)
	suite.NoError(err)

	suite.Len(respData.Apps, 1)
	app_1 := respData.Apps[0]
	suite.Equal("app-a", app_1.Name)

	appSvc.AssertExpectations(suite.T())
}

func (suite *AppTestSuite) Test_HandleRestartApp_Ok() {
	userName := "test-user"
	token, err := auth.Login(userName)
	suite.NoError(err)

	appSvc := new(mocks.AppSvc)
	projectSvc := new(mocks.ProjectSvc)
	root := &handlers.Root{AppSvc: appSvc, ProjectSvc: projectSvc}

	appName := "app-a"

	projName := "project-a"
	proj := &responses.Project{
		Name: projName,
	}

	projectSvc.On("Get", mock.AnythingOfType("*context.valueCtx"), userName, projName).Return(proj, nil)
	appSvc.On("Restart", mock.AnythingOfType("*context.valueCtx"), projName, appName).Return(nil)

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	req, err := http.NewRequest(http.MethodPost, s.URL+fmt.Sprintf("/v1/projects/%s/apps/%s/restart", projName, appName), nil)
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

	appSvc.AssertExpectations(suite.T())
}

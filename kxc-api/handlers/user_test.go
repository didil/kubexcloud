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
	"github.com/didil/kubexcloud/kxc-api/services"
	"github.com/didil/kubexcloud/kxc-api/testsupport"
	"github.com/didil/kubexcloud/kxc-api/testsupport/auth"
	"github.com/didil/kubexcloud/kxc-api/testsupport/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
}

func (suite *UserTestSuite) SetupSuite() {
	testsupport.BootstrapTests("../.env.test")
}
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) Test_HandleLoginUser_Ok() {
	userSvc := new(mocks.UserSvc)
	root := &handlers.Root{UserSvc: userSvc}

	reqData := &requests.LoginUser{
		Name:     "test-user",
		Password: "123456",
	}

	token := "TEST_AUTH_TOKEN"

	userSvc.On("Login", mock.AnythingOfType("*context.valueCtx"), reqData.Name, reqData.Password).Return(token, nil)

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, s.URL+"/v1/users/login", &b)
	suite.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)

	defer resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))

	var respData *responses.LoginUser
	err = json.NewDecoder(resp.Body).Decode(&respData)
	suite.NoError(err)

	suite.Equal(token, respData.Token)

	userSvc.AssertExpectations(suite.T())
}

func (suite *UserTestSuite) Test_HandleCreateUser_Ok() {
	userName := "adminUser"

	token, err := auth.Login(userName)
	suite.NoError(err)

	userSvc := new(mocks.UserSvc)
	userSvc.On("HasRole", mock.AnythingOfType("*context.valueCtx"), userName, services.UserRoleAdmin).Return(true, nil)

	root := &handlers.Root{UserSvc: userSvc}

	reqData := &requests.CreateUser{
		Name:     "test-user",
		Password: "123456",
		Role:     services.UserRoleRegular,
	}

	userSvc.On("Create", mock.AnythingOfType("*context.valueCtx"), reqData).Return(nil)

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, s.URL+"/v1/users", &b)
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

	userSvc.AssertExpectations(suite.T())
}

func (suite *UserTestSuite) Test_HandleCreateUser_NotAdmin() {
	userName := "adminUser"

	token, err := auth.Login(userName)
	suite.NoError(err)

	userSvc := new(mocks.UserSvc)
	userSvc.On("HasRole", mock.AnythingOfType("*context.valueCtx"), userName, services.UserRoleAdmin).Return(false, nil)

	root := &handlers.Root{UserSvc: userSvc}

	reqData := &requests.CreateUser{
		Name:     "test-user",
		Password: "123456",
		Role:     services.UserRoleRegular,
	}

	r := api.BuildRouter(root)
	s := httptest.NewServer(r)
	defer s.Close()

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, s.URL+"/v1/users", &b)
	suite.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)

	defer resp.Body.Close()
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)

	userSvc.AssertExpectations(suite.T())
}

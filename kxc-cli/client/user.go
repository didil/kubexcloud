package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/didil/kubexcloud/kxc-api/requests"
	"github.com/didil/kubexcloud/kxc-api/responses"
)

func (cl *Client) LoginUser(apiURL, userName, password string) (string, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return "", fmt.Errorf("invalid api url %v : %v", apiURL, err)
	}

	u.Path = path.Join(u.Path, "v1/users/login")

	reqData := &requests.LoginUser{
		Name:     userName,
		Password: password,
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(reqData)
	if err != nil {
		return "", fmt.Errorf("encode req data: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), &b)
	if err != nil {
		return "", fmt.Errorf("new req: %v", err)
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("req do: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		errData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("http read: %v", err)
		}

		return "", fmt.Errorf("http: %v, %s", resp.StatusCode, string(errData))
	}

	respData := &responses.LoginUser{}

	err = json.NewDecoder(resp.Body).Decode(respData)
	if err != nil {
		return "", fmt.Errorf("decode: %v", err)
	}

	return respData.Token, nil
}

func (cl *Client) CreateUser(userName, password, role string) error {
	u, err := url.Parse(cl.apiURL)
	if err != nil {
		return fmt.Errorf("invalid api url %v : %v", cl.apiURL, err)
	}

	u.Path = path.Join(u.Path, "v1/users")

	reqData := &requests.CreateUser{
		Name:     userName,
		Password: password,
		Role:     role,
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(reqData)
	if err != nil {
		return fmt.Errorf("encode req data: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), &b)
	if err != nil {
		return fmt.Errorf("new req: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+cl.authToken)

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("req do: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		errData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("http read: %v", err)
		}

		return fmt.Errorf("http: %v, %s", resp.StatusCode, string(errData))
	}

	return nil
}

func (cl *Client) ListUsers() (*responses.ListUser, error) {
	u, err := url.Parse(cl.apiURL)
	if err != nil {
		return nil, fmt.Errorf("invalid api url %v : %v", cl.apiURL, err)
	}

	u.Path = path.Join(u.Path, "v1/users")

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new req: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+cl.authToken)

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("req do: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		errData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("http read: %v", err)
		}

		return nil, fmt.Errorf("http: %v, %s", resp.StatusCode, string(errData))
	}

	respData := &responses.ListUser{}

	err = json.NewDecoder(resp.Body).Decode(respData)
	if err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	return respData, nil
}

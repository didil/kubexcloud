package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/didil/kubexcloud/kxc-api/responses"
)

func (cl *Client) ListApps(projectName string) (*responses.ListApp, error) {
	u, err := url.Parse(cl.apiURL)
	if err != nil {
		return nil, fmt.Errorf("invalid api url %v : %v", cl.apiURL, err)
	}

	u.Path = path.Join(u.Path, fmt.Sprintf("v1/projects/%s/apps", projectName))

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

	respData := &responses.ListApp{}

	err = json.NewDecoder(resp.Body).Decode(respData)
	if err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	return respData, nil
}

func (cl *Client) RestartApp(projectName, appName string) error {
	u, err := url.Parse(cl.apiURL)
	if err != nil {
		return fmt.Errorf("invalid api url %v : %v", cl.apiURL, err)
	}

	u.Path = path.Join(u.Path, fmt.Sprintf("v1/projects/%s/apps/%s/restart", projectName, appName))

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
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

	respData := &responses.ListApp{}

	err = json.NewDecoder(resp.Body).Decode(respData)
	if err != nil {
		return fmt.Errorf("decode: %v", err)
	}

	return nil
}

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

func (cl *Client) ListProjects() (*responses.ListProject, error) {
	u, err := url.Parse(cl.apiURL)
	if err != nil {
		return nil, fmt.Errorf("invalid api url %v : %v", cl.apiURL, err)
	}

	u.Path = path.Join(u.Path, "v1/projects")

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

	respData := &responses.ListProject{}

	err = json.NewDecoder(resp.Body).Decode(respData)
	if err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	return respData, nil
}

package proxy

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/alexellis/jaas/pkg/types"
)

func NewJaaSClient(config types.JaaSConfig) JaaSClient {
	client := http.Client{}
	return JaaSClient{Client: &client, Config: config}
}

type JaaSClient struct {
	Client *http.Client
	Config types.JaaSConfig
}

func (j *JaaSClient) List(address string) ([]types.Task, error) {
	tasks := []types.Task{}

	if len(address) == 0 {
		return tasks, fmt.Errorf("bad address")
	}

	u, _ := url.Parse(address)
	u.Path = path.Join(u.Path, "list")

	req, reqErr := http.NewRequest(http.MethodGet, u.String(), nil)

	if reqErr != nil {
		return tasks, reqErr
	}

	for _, auth := range j.Config.Auths {
		if auth.Address == address {
			p, _ := hex.DecodeString(auth.Password)
			req.SetBasicAuth(auth.Username, string(p))
			break
		}
	}

	res, err := j.Client.Do(req)
	if err != nil {
		return tasks, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	bodyBytes, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(bodyBytes))
	}

	marshalErr := json.Unmarshal(bodyBytes, &tasks)
	if marshalErr != nil {
		return tasks, marshalErr
	}

	return tasks, nil
}

func (j *JaaSClient) Run(address string, runReq RunRequest) (int, error) {

	if len(address) == 0 {
		return http.StatusBadGateway, fmt.Errorf("bad address")
	}

	u, _ := url.Parse(address)
	u.Path = path.Join(u.Path, "run")

	reqBytes, marshalErr := json.Marshal(runReq)
	if marshalErr != nil {
		return http.StatusBadRequest, marshalErr
	}

	req, reqErr := http.NewRequest(http.MethodGet, u.String(), bytes.NewReader(reqBytes))
	if reqErr != nil {
		return http.StatusBadGateway, reqErr
	}

	for _, auth := range j.Config.Auths {
		if auth.Address == address {
			p, _ := hex.DecodeString(auth.Password)
			req.SetBasicAuth(auth.Username, string(p))
			break
		}
	}

	res, err := j.Client.Do(req)
	if err != nil {
		return http.StatusBadGateway, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	return res.StatusCode, nil
}

func (j *JaaSClient) Ping(address string) (*PingResponse, error) {
	if len(address) == 0 {
		return nil, fmt.Errorf("bad address")
	}
	u, _ := url.Parse(address)
	u.Path = path.Join(u.Path, "ping")

	req, reqErr := http.NewRequest(http.MethodGet, u.String(), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	for _, auth := range j.Config.Auths {
		if auth.Address == address {
			p, _ := hex.DecodeString(auth.Password)
			req.SetBasicAuth(auth.Username, string(p))
			break
		}
	}

	res, err := j.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	return &PingResponse{Status: res.StatusCode}, nil
}

type PingResponse struct {
	Status int
	Error  error
}

type RunRequest struct {
	Job *JaasJob
}

type JaasJob struct {
	Image   string
	Command string
}

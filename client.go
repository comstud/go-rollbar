package rollbar

import (
	"bytes"
	"encoding/json"
	"net/http"
	net_url "net/url"
)

const DEFAULT_NOTIFIER_NAME = "go-rollbar"
const DEFAULT_NOTIFIER_VERSION = "0.0.1"
const DEFAULT_NOTIFIER_ENVIRONMENT = "production"

const DEFAULT_API_BASE_URL = "https://api.rollbar.com/api/1"

type Client struct {
	apiBaseURL  string
	accessToken string

	notifierName    string
	notifierVersion string
	environment     string
	codeVersion     string

	httpClient *http.Client
}

func (self *Client) Environment() string {
	return self.environment
}

func (self *Client) SetEnvironment(env string) *Client {
	self.environment = env
	return self
}

func (self *Client) CodeVersion() string {
	return self.codeVersion
}

func (self *Client) SetCodeVersion(code_version string) *Client {
	self.codeVersion = code_version
	return self
}

func (self *Client) APIBaseURL() string {
	return self.apiBaseURL
}

func (self *Client) SetAPIBaseURL(base_url string) *Client {
	self.apiBaseURL = base_url
	return self
}

func (self *Client) httpGet(url string, query net_url.Values, resp interface{}) error {
	if query == nil {
		query = make(net_url.Values)
	}

	query.Set("access_token", self.accessToken)

	if url != "" && url[0] != '/' {
		url = self.apiBaseURL + "/" + url
	} else {
		url = self.apiBaseURL + url
	}

	query_str := query.Encode()
	if query_str != "" {
		url += "?" + query_str
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")

	http_resp, err := self.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer http_resp.Body.Close()

	return json.NewDecoder(http_resp.Body).Decode(resp)
}

func (self *Client) httpPatch(url string, data interface{}, resp interface{}) error {
	query := net_url.Values{
		"access_token": []string{self.accessToken},
	}

	if url != "" && url[0] != '/' {
		url = self.apiBaseURL + "/" + url
	} else {
		url = self.apiBaseURL + url
	}

	query_str := query.Encode()
	if query_str != "" {
		url += "?" + query_str
	}

	json_data, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	http_resp, err := self.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer http_resp.Body.Close()
	return json.NewDecoder(http_resp.Body).Decode(resp)
}

func NewClient(access_token string) (*Client, error) {
	return &Client{
		apiBaseURL:      DEFAULT_API_BASE_URL,
		accessToken:     access_token,
		notifierName:    DEFAULT_NOTIFIER_NAME,
		notifierVersion: DEFAULT_NOTIFIER_VERSION,
		environment:     DEFAULT_NOTIFIER_ENVIRONMENT,
		httpClient:      &http.Client{},
	}, nil
}

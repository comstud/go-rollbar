package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	net_url "net/url"
)

const (
	DEFAULT_API_BASE_URL         = "https://api.rollbar.com/api/1"
	DEFAULT_NOTIFIER_NAME        = "go-rollbar"
	DEFAULT_NOTIFIER_VERSION     = "0.0.1"
	DEFAULT_NOTIFIER_ENVIRONMENT = "development"
)

// Client struct
type Client struct {
	apiBaseURL  string
	accessToken string

	notifierName    string
	notifierVersion string
	environment     string
	codeVersion     string

	httpClient *http.Client
}

func (self *Client) decodeBody(http_resp *http.Response, resp interface{}) error {
	if http_resp.StatusCode >= 200 && http_resp.StatusCode < 300 {
		return json.NewDecoder(http_resp.Body).Decode(resp)
	}

	// Hrmph.
	data, err := ioutil.ReadAll(http_resp.Body)
	if err == nil {
		return fmt.Errorf("Got code %d: %s\n", http_resp.StatusCode, string(data))
	} else {
		return fmt.Errorf("Got code %d: <Error reading body: %s>\n", http_resp.StatusCode, err)
	}
}

func (self *Client) httpCall(method string, url string, query net_url.Values, data interface{}, resp interface{}) error {
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

	var json_data []byte

	if data != nil {
		var err error

		json_data, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	http_resp, err := self.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer http_resp.Body.Close()

	return self.decodeBody(http_resp, resp)

}

func (self *Client) httpGet(url string, query net_url.Values, resp interface{}) error {
	return self.httpCall("GET", url, query, nil, resp)
}

func (self *Client) httpPatch(url string, data interface{}, resp interface{}) error {
	return self.httpCall("PATCH", url, nil, data, resp)
}

func (self *Client) httpPost(url string, data interface{}, resp interface{}) error {
	return self.httpCall("POST", url, nil, data, resp)
}

// Get the environment
func (self *Client) Environment() string {
	return self.environment
}

// Set the environment
func (self *Client) SetEnvironment(env string) *Client {
	self.environment = env
	return self
}

// Get the code version
func (self *Client) CodeVersion() string {
	return self.codeVersion
}

// Set the code version
func (self *Client) SetCodeVersion(code_version string) *Client {
	self.codeVersion = code_version
	return self
}

// Get the base API URL
func (self *Client) APIBaseURL() string {
	return self.apiBaseURL
}

// Set the base API URL
func (self *Client) SetAPIBaseURL(base_url string) *Client {
	self.apiBaseURL = base_url
	return self
}

// Create a new client with specified access token
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

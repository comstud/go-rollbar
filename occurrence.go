package rollbar

import (
	"errors"
	"fmt"
	"net/url"
)

type OccurrenceMetadata struct {
	CustomerTimestamp JSONTime    `json:"customer_timestamp"`
	TimestampMS       int64       `json:"timestamp_ms"`
	APIServerHostname string      `json:"api_server_hostname"`
	Debug             interface{} `json:"debug"`
}

type OccurrenceNotifier struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type OccurrenceBody struct {
	Message map[string]interface{} `json:"message"`
}

type OccurrenceData struct {
	Environment string             `json:"environment"`
	Timestamp   JSONTime           `json:"timestamp"`
	CodeVersion string             `json:"code_version"`
	Platform    string             `json:"platform"`
	Level       string             `json:"level"`
	Notifier    OccurrenceNotifier `json:"notifier"`
	Context     string             `json:"context"`
	Title       string             `json"title"`
	Body        OccurrenceBody     `json:"body"`
	Metadata    OccurrenceMetadata `json:"metadata"`
	Framework   string             `json:"framework"`
	UUID        string             `json:"uuid"`
}

type Occurrence struct {
	ID         uint64         `json:"id"`
	Project_id uint64         `json:"project_id"`
	Timestamp  JSONTime       `json"timestamp"`
	Version    uint           `json:"version"`
	Data       OccurrenceData `json:"data"`
	Billable   uint           `json:"billable"`
}

func (self *Occurrence) String() string {
	return self.AsPrettyJSON()
}

func (self *Occurrence) AsJSON() string {
	return asJSON(self)
}

func (self *Occurrence) AsPrettyJSON() string {
	return asPrettyJSON(self)
}

type OccurrenceResponse struct {
	BaseAPIResponse
	*Occurrence `json:"result"`
}

type OccurrencesResult struct {
	Occurrences []*Occurrence `json:"instances"`
}

type OccurrencesResponse struct {
	rollbar *Client
	BaseAPIResponse
	Page               uint
	*OccurrencesResult `json:"result"`
}

func (self *OccurrencesResponse) String() string {
	return self.AsPrettyJSON()
}

func (self *OccurrencesResponse) AsJSON() string {
	s := "["
	for i, occur := range self.Occurrences {
		if i != 0 {
			s += ","
		}
		s += asJSON(occur)
	}
	return s + "]"
}

func (self *OccurrencesResponse) AsPrettyJSON() string {
	s := "["
	for i, occur := range self.Occurrences {
		if i != 0 {
			s += ","
		}
		s += asPrettyJSON(occur)
	}
	return s + "]"
}

func (self *OccurrencesResponse) HasMorePages() bool {
	return self.IsSuccess() && len(self.Occurrences) > 0
}

func (self *OccurrencesResponse) GetNextPage() (*OccurrencesResponse, error) {
	if !self.HasMorePages() {
		return self, nil
	}
	resp := &OccurrencesResponse{
		rollbar: self.rollbar,
		Page:    self.Page + 1,
	}
	return self.rollbar.getOccurrences(resp)
}

func (self *Client) GetOccurrence(id uint64) (*OccurrenceResponse, error) {
	occur_resp := &OccurrenceResponse{}

	err := self.httpGet(
		fmt.Sprintf("/instance/%d", id),
		nil,
		&occur_resp,
	)
	if err != nil {
		return nil, err
	}

	return occur_resp, nil
}

func (self *Client) getOccurrences(resp *OccurrencesResponse) (*OccurrencesResponse, error) {
	query := url.Values{"page": []string{fmt.Sprintf("%d", resp.Page)}}

	err := self.httpGet("/instances", query, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (self *Client) GetOccurrences() (*OccurrencesResponse, error) {
	resp := &OccurrencesResponse{
		rollbar: self,
		Page:    1,
	}
	return self.getOccurrences(resp)
}

func (self *Client) GetOccurrencesWithPage(page uint) (*OccurrencesResponse, error) {
	if page == 0 {
		return nil, errors.New("Page must be greater than 0")
	}
	resp := &OccurrencesResponse{
		rollbar: self,
		Page:    page,
	}
	return self.getOccurrences(resp)
}

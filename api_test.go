package main

import (
	"net/http"
	"testing"
)

// A mock client for functions relying on mocked values of bitlyClient interface.
// Tests setup desired response of sendRequest by setting the
// sendRequestResponse and err field.  They can also simulate paginated
// responses by setting the Pages field to an int greater than 1.
type mockBitlyClient struct {
	sendRequestResponse []byte
	err                 error
	pages               int
}

// createRequest can return mocked values but we are already testing
// that error handling scenario so it suffices to just return a real request.
func (c *mockBitlyClient) createRequest(path, verb, body string) (*http.Request, error) {
	return http.NewRequest(verb, path, nil)
}

// sendRequest returns the response that the test sets up or a paginated version
// of the response.
func (c *mockBitlyClient) sendRequest(req *http.Request) ([]byte, error) {

	// for pagination send links containing next urls to verify that they are followed
	if c.pages > 1 {
		c.pages--
		return []byte(`{"pagination": {"next": "https://api-ssl.bitly.com/v4/groups/Bk1hmw/bitlinks"},
				 "links": [{"link": "http://bit.ly/HGDGAX","id": "bit.ly/HGDGAX"}]}`), nil
	}
	return c.sendRequestResponse, c.err
}

// Deserialization tests for response structs
func Test_DeserializeUserResponse(t *testing.T) {

	response := []byte(`{"name":"carrpet912","default_group_guid":"Bk1hmwBHfQK"}`)
	toTest := &UserInfo{}
	err := toTest.deserialize(response)
	if err != nil {
		t.Error(err)
	}
	if toTest.GroupGuid != "Bk1hmwBHfQK" {
		t.Error(`UserInfo.deserialize failed`)
	}
}

func Test_DeserializeBitlinksByGroupResponse(t *testing.T) {
	// assert that we are getting and storing the link's id
	response := []byte(`{"pagination": {
      "total":1,
      "page":0,
      "next": "http://next.com"
    }, "links": [{"link": "http://bit.ly/HGDGAX","id": "1928432"}]}`)
	toTest := &GroupBitlinks{}
	err := toTest.deserialize(response)
	if err != nil {
		t.Error(err)
	}
	if toTest.Links[0].Link != "http://bit.ly/HGDGAX" {
		t.Error(`Bitlink.deserialize failed: links`)
	}
	if toTest.Links[0].ID != "1928432" {
		t.Error(`Bitlink.deserialize failed: ids`)
	}
	if toTest.Pagination.Next != "http://next.com" {
		t.Error(`Bitlink.deserialize failed: pagination`)
	}
}

func Test_DeserializeClickMetrics(t *testing.T) {
	response := []byte(
		`{"units":30, "metrics": [{"clicks":27,"value": "US"}, {"clicks": 1000, "value": "China"}]}`)
	toTest := &ClickMetrics{}
	err := toTest.deserialize(response)
	if err != nil {
		t.Error(err)
	}
	if len(toTest.Metrics) != 2 {
		t.Error(`ClickMetrics.deserialize failed: metrics length`)
	}
	if toTest.Metrics[0].Clicks != 27 {
		t.Error(`ClickMetrics.deserialize failed: clicks`)

	}
	if toTest.Metrics[0].Country != "US" {
		t.Error(`ClickMetrics.deserialize failed: country`)
	}
}

// Tests for Bitly API methods
func Test_GetUserInfo(t *testing.T) {
	mock := &mockBitlyClient{}
	api := bitlinksMetricsAPI{}
	api.GetUserInfo(mock)
}

func Test_GetBitlinksForGroupSinglePage(t *testing.T) {
	resp := []byte(`{"links": [{"link": "http://bit.ly/UFHISO","id": "bit.ly/UFHISO"},
		{"link": "http://nyti.ms/2GnOpXm","id": "nyti.ms/2GnOpXm"}],
		"pagination": {
          "total":2,
          "page":1,
					"size":50,
          "next":"",
					"prev":""}}`)
	mock := &mockBitlyClient{sendRequestResponse: resp}
	api := bitlinksMetricsAPI{}
	bl, err := api.GetBitlinksForGroup(mock, "ABC3DgEF")
	if err != nil {
		t.Errorf("GetBitlinksForGroup: expected success, received error: %s", err.Error())
	}
	if len(bl.Links) != 2 {
		t.Errorf("GetBitlinksForGroup: return value expected %d links, received %d", 2, len(bl.Links))
	}
	if bl.Links[0].Link != "http://bit.ly/UFHISO" || bl.Links[1].Link != "http://nyti.ms/2GnOpXm" {
		t.Errorf("GetBitlinksForGroup: return value wrong links value, received link %s and %s", bl.Links[0], bl.Links[1])
	}
}

func Test_GetBitlinksByGroupPagination(t *testing.T) {
	singlePageResp := []byte(`{"links": [{"link": "http://bit.ly/UFHISO","id": "bit.ly/UFHISO"},
		{"link": "http://nyti.ms/2GnOpXm","id": "nyti.ms/2GnOpXm"}],
		"pagination": {
          "total":6,
          "page":1,
					"size":50,
          "next":"",
					"prev":""}}`)
	mock := &mockBitlyClient{pages: 4, sendRequestResponse: singlePageResp}
	api := bitlinksMetricsAPI{}
	toTest, err := api.GetBitlinksForGroup(mock, "ABC3DgEF")
	if err != nil {
		t.Error(err)
	}
	// the last page returns two links and each paginated response has one link
	// for a total of 5
	if len(toTest.Links) != 5 {
		t.Errorf(`GetBitlinksByGroupPagination failed: number of links, expected: %d received: %d`, 5, len(toTest.Links))
	}
}

func Test_GetBitlinkClicksByCountry(t *testing.T) {
	resp := []byte(`{"units":30, "facet":"countries", "unit":"day",
		"metrics": [{"clicks":20,"value":"US"}, {"clicks":809,"value":"Mexico"}]}`)
	mock := &mockBitlyClient{sendRequestResponse: resp}
	api := bitlinksMetricsAPI{}
	toTest, err := api.GetBitlinkClicksByCountry(mock, Bitlink{})
	if err != nil {
		t.Errorf("GetBitlinkClicksByCountry: expected success, received error: %s", err.Error())
	}
	if len(toTest.Metrics) != 2 {
		t.Errorf("GetBitlinkClicksByCountry: expected metrics length 2, received length: %d", len(toTest.Metrics))
	}
}

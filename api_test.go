package main

import (
	"net/http"
	"testing"
)

type mockBitlyClient struct {
	Pages int
}

func (c *mockBitlyClient) createRequest(path, verb, body string) (*http.Request, error) {
	return http.NewRequest(verb, path, nil)

}
func (c *mockBitlyClient) sendRequest(req *http.Request) ([]byte, error) {
	if req.URL.Path == "/v4/groups/"+"Bk1hmwBHfQK"+"/bitlinks" {
		// send a paginated link, else send a single link
		if c.Pages > 0 {
			c.Pages--
			return []byte(`{"pagination": {"next": "https://api-ssl.bitly.com/v4/groups/Bk1hmwBHfQK/bitlinks"},
				 "links": [{"link": "http://bit.ly/HGDGAX","id": "bit.ly/HGDGAX"}]}`), nil
		}
		return []byte(`{"pagination": {
          "total":1,
          "page":0,
          "next":""
        }, "links": [{"link": "http://bit.ly/UFHISO","id": "bit.ly/UFHISO"}]}`), nil
	}

	return []byte{}, nil

}

func Test_GetUserInfo(t *testing.T) {
	mock := &mockBitlyClient{}
	api := bitlinksMetricsAPI{}
	api.GetUserInfo(mock)
}

// test that deserialization works
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

func Test_GetBitlinksByGroupPagination(t *testing.T) {
	mock := &mockBitlyClient{Pages: 3}
	api := bitlinksMetricsAPI{}
	toTest, err := api.GetBitlinksForGroup(mock, "Bk1hmwBHfQK")
	if err != nil {
		t.Error(err)
	}
	if len(toTest.Links) != 4 {
		t.Error(`GetBitlinksByGroupPagination failed`)
	}
}

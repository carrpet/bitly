package client

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

// test that deserialization works
func Test_DeserializeUserResponse(t *testing.T) {

	response := []byte(`{"name":"carrpet912","default_group_guid":"Bk1hmwBHfQK"}`)
	toTest := &bitlyUserInfo{}
	err := toTest.deserialize(response)
	if err != nil {
		t.Error(err)
	}
	if toTest.GroupGuid != "Bk1hmwBHfQK" {
		t.Error(`bitlyUserInfo.deserialize failed`)
	}

}

func Test_DeserializeBitlinksByGroupResponse(t *testing.T) {
	// assert that we are getting and storing the link's id
	response := []byte(`{"pagination": {
      "total":1,
      "page":0,
      "next": "http://next.com"
    }, "links": [{"link": "http://bit.ly/HGDGAX","id": "1928432"}]}`)
	toTest := &bitlyGroupsBitLinks{}
	err := toTest.deserialize(response)
	if err != nil {
		t.Error(err)
	}
	if toTest.Links[0].Link != "http://bit.ly/HGDGAX" {
		t.Error(`bitlyBitlinks.deserialize failed: links`)
	}
	if toTest.Links[0].ID != "1928432" {
		t.Error(`bitlyBitlinks.deserialize failed: ids`)
	}
	if toTest.Pagination.Next != "http://next.com" {
		t.Error(`bitlyBitlinks.deserialize failed: pagination`)

	}

}

func Test_GetBitlinksByGroupPagination(t *testing.T) {
	mock := &mockBitlyClient{Pages: 3}
	toTest, err := GetBitlinksForGroup(mock, "Bk1hmwBHfQK")
	if err != nil {
		t.Error(err)
	}
	if len(toTest.Links) != 4 {
		t.Error(`GetBitlinksByGroupPagination failed`)
	}
}

func Test_GetUserInfo(t *testing.T) {
	mock := &mockBitlyClient{}
	GetUserInfo(mock)
}

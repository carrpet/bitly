package client

import (
	"net/http"
	"testing"
)

type mockBitlyClient struct{}

func (c mockBitlyClient) createRequest(path, verb, body string) (*http.Request, error) {
	return &http.Request{}, nil

}
func (c mockBitlyClient) sendRequest(req *http.Request) ([]byte, error) {
	return []byte{}, nil

}

/*
func generateRandomLinksResponse(length int) []byte {
  res := []byte
  for i in length {
    res.append(`{"link": "http something", "id": "somerandomstring" `)
  }
  return res
}
*/

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
      "page":0
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

}

func Test_GetUserInfo(t *testing.T) {
	mock := mockBitlyClient{}
	GetUserInfo(mock)
}

func Test_GetBitlinksForGroup(t *testing.T) {
	mock := mockBitlyClient{}
	GetBitlinksForGroup(mock, "1928432")
}

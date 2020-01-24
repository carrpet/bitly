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

// test that deserialization works
func Test_DeserializeResponse(t *testing.T) {

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

func Test_GetUserInfo(t *testing.T) {
	mock := mockBitlyClient{}
	GetUserInfo(mock)

}

// test that we extract the group guid from the response
/*
func Test_GetJsonResponse(t *testing.T) {
	client := &mockBitlyClient{}
	client.GetUserInfo(client)

}*/

package client

import (
	"testing"
)

// test that deserialization works
func Test_DeserializeResponse(t *testing.T) {

	response := []byte(`{"name":"carrpet912","default_group_guid":"Bk1hmwBHfQK"}`)
	toTest := &BitlyUserInfo{}
	err := toTest.Deserialize(response)
	if err != nil {
		t.Error(err)
	}
	if toTest.GroupGuid != "Bk1hmwBHfQK" {
		t.Error(`BitlyUserInfo.Deserialize failed`)
	}

}

// test that we extract the group guid from the response
/*
func Test_GetJsonResponse(t *testing.T) {
	client := &mockBitlyClient{}
	client.GetUserInfo(client)

}*/

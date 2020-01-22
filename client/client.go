package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type BitlyObject interface {
	Deserialize(res []byte)
}

type BitlyClientRequest struct {
	req *http.Request
}

func newBitlyClientRequest(path, verb, token string) (*BitlyClientRequest, error) {
	baseURL := "https://api-ssl.bitly.com/v4/"
	req, err := http.NewRequest(verb, baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)
	return &BitlyClientRequest{req: req}, nil
}

type BitlyUserInfo struct {
	GroupGuid string `json:"default_group_guid"`
	Name      string `json:"name"`
}

func (o *BitlyUserInfo) Deserialize(res []byte) error {

	if err := json.Unmarshal(res, &o); err != nil {
		return err
	}
	return nil

}

func GetUserInfo(token string) (*BitlyUserInfo, error) {
	req, err := newBitlyClientRequest("user", "GET", token)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req.req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	userInfo := &BitlyUserInfo{}
	err = userInfo.Deserialize(body)
	if err != nil {
		return nil, err
	}
	return userInfo, nil

}

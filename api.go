package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const bitlyAPI = "https://api-ssl.bitly.com/v4/"

type BitlinksMetricsAPI struct{}

type bitlyClient interface {
	createRequest(path, verb, body string) (*http.Request, error)
	// SendRequest sends the request to the server and
	// reads the full response into a byte array
	sendRequest(*http.Request) ([]byte, error)
}

type BitlinksMetrics interface {
	GetUserInfo(bitlyClient) (*bitlyUserInfo, error)
	GetBitlinksForGroup(bitlyClient, string) (*bitlyGroupsBitLinks, error)
	GetClicksByCountry(bitlyClient, Bitlink) (*ClickMetrics, error)
}

type BitlyClientInfo struct {
	Token string
}

type bitlyObject interface {
	deserialize(res []byte)
}

type bitlyUserInfo struct {
	GroupGuid string `json:"default_group_guid"`
	Name      string `json:"name"`
}

type bitlyGroupsBitLinks struct {
	Pagination Pagination `json:"pagination"`
	Links      []Bitlink  `json:"links"`
}
type Bitlink struct {
	Link string `json:"link"`
	ID   string `json:"id"`
}

type Pagination struct {
	Next  string `json:"next"`
	Total int    `json:"total"`
}

type ClickMetrics struct {
	//units by default are in days, ie the time range
	Units   int            `json:"units"`
	Metrics []CountryClick `json:"metrics"`
}

type CountryClick struct {
	Clicks  int    `json:"clicks"`
	Country string `json:"value"`
}

// all package exported methods
func (b *BitlinksMetricsAPI) GetUserInfo(client bitlyClient) (*bitlyUserInfo, error) {
	req, err := client.createRequest(bitlyAPI+"user", "GET", "")
	if err != nil {
		return nil, err
	}
	body, err := client.sendRequest(req)
	if err != nil {
		return nil, err
	}
	userInfo := &bitlyUserInfo{}
	err = userInfo.deserialize(body)
	if err != nil {
		return nil, err
	}
	return userInfo, nil

}

func (b *BitlinksMetricsAPI) GetBitlinksForGroup(client bitlyClient, groupGUID string) (*bitlyGroupsBitLinks, error) {
	path := bitlyAPI + "groups/" + groupGUID + "/bitlinks"
	verb := "GET"
	req, err := client.createRequest(path, verb, "")
	if err != nil {
		return nil, err
	}
	body, err := client.sendRequest(req)
	if err != nil {
		return nil, err
	}
	bitlinks := &bitlyGroupsBitLinks{}
	err = bitlinks.deserialize(body)
	if err != nil {
		fmt.Printf("error is: %s", err.Error())
		return nil, err
	}

	links := append([]Bitlink{}, bitlinks.Links...)
	for bitlinks.Pagination.Next != "" {
		req, err = client.createRequest(bitlinks.Pagination.Next, verb, "")
		if err != nil {
			return nil, err
		}
		body, err = client.sendRequest(req)
		if err != nil {
			return nil, err
		}
		bitlinks = &bitlyGroupsBitLinks{}
		err = bitlinks.deserialize(body)
		if err != nil {
			return nil, err
		}

		links = append(links, bitlinks.Links...)
	}

	bitlinks.Links = links
	return bitlinks, nil

}

func (b *BitlinksMetricsAPI) GetClicksByCountry(client bitlyClient, link Bitlink) (*ClickMetrics, error) {
	path := bitlyAPI + "bitlinks/" + link.ID + "/countries"
	params := "?units=30"
	verb := "GET"
	req, err := client.createRequest(path+params, verb, "")
	if err != nil {
		return nil, err
	}
	body, err := client.sendRequest(req)
	if err != nil {
		return nil, err
	}
	metrics := &ClickMetrics{}
	err = metrics.deserialize(body)
	if err != nil {
		fmt.Printf("error is: %s", err.Error())
		return nil, err
	}
	return metrics, nil

}

func (o *bitlyUserInfo) deserialize(res []byte) error {

	if err := json.Unmarshal(res, &o); err != nil {
		return err
	}
	return nil

}
func (o *bitlyGroupsBitLinks) deserialize(res []byte) error {

	if err := json.Unmarshal(res, &o); err != nil {
		return err
	}
	return nil

}

func (o *ClickMetrics) deserialize(res []byte) error {

	if err := json.Unmarshal(res, &o); err != nil {
		return err
	}
	return nil

}

//all package internal methods
func (c *BitlyClientInfo) createRequest(path, verb, body string) (*http.Request, error) {
	req, err := http.NewRequest(verb, path, nil)
	if err != nil {
		return nil, err
	}
	bearer := "Bearer " + c.Token
	req.Header.Add("Authorization", bearer)
	return req, nil
}

func (c *BitlyClientInfo) sendRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	return ioutil.ReadAll(resp.Body)
}

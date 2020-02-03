package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockBitlinksMetricsAPI struct {
	getUserInfo              *UserInfo
	getUserInfoError         error
	getBitlinksForGroup      *GroupBitlinks
	getBitlinksForGroupError error
	getClicksByCountry       *ClickMetrics
	getClicksByCountryError  error
}

func (m *mockBitlinksMetricsAPI) GetUserInfo(c bitlyClient) (*UserInfo, error) {
	if m.getUserInfoError != nil {
		return nil, m.getUserInfoError
	}
	return &UserInfo{GroupGuid: "abcdefgh", Name: "petertest"}, nil

}

func (m *mockBitlinksMetricsAPI) GetBitlinksForGroup(c bitlyClient, guid string) (*GroupBitlinks, error) {
	if m.getBitlinksForGroupError != nil {
		return nil, m.getBitlinksForGroupError
	}
	l := []Bitlink{{Link: "http://something", ID: "something"},
		{Link: "http://something1", ID: "something1"},
		{Link: "http://something2", ID: "something2"},
	}
	return &GroupBitlinks{Links: l}, nil

}

func (m *mockBitlinksMetricsAPI) GetBitlinkClicksByCountry(c bitlyClient, b Bitlink) (*ClickMetrics, error) {
	cc := []CountryClick{
		{Clicks: 80, Country: "US"},
		{Clicks: 40, Country: "Argentina"},
		{Clicks: 10, Country: "Norway"},
		{Clicks: 60, Country: "Sweden"},
		{Clicks: 1200, Country: "China"},
	}
	return &ClickMetrics{Units: 30, Metrics: cc}, nil

}

// Main handler tests

func TestAvgClickMetricsHandler(t *testing.T) {

	expected := map[string]int{
		"US":        80 * 3,
		"Argentina": 40 * 3,
		"Norway":    10 * 3,
		"Sweden":    60 * 3,
		"China":     1200 * 3,
	}
	context := BitlyClientInfo{}
	mock := &mockBitlinksMetricsAPI{}

	cc, err := context.avgClicks(mock)
	if err != nil {
		t.Errorf("avgClicks returned error: %s", err.Error())
	}
	if len(cc) != 5 {
		t.Errorf("avgClicks return value has length %d, expected length 5", len(cc))
	}
	for _, item := range cc {
		if _, ok := expected[item.Country]; !ok {
			t.Errorf("avgClicks error: Country %s not found in returned results", item.Country)
		}

		if expected[item.Country] != item.Clicks {
			t.Errorf("avgClicks error: Number of clicks for country %s: expected %d, received %d", item.Country, expected[item.Country], item.Clicks)
		}
	}
}

// HTTP Server tests
func TestHandleAvgMetrics(t *testing.T) {
	req, err := http.NewRequest("GET", "/groups/Bk1hmwBHfQK/countries/averages", nil)
	req.Header.Add("Authorization", "Bearer myfakeaccesstoken")
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	c := BitlyClientInfo{}
	api := &bitlinksMetricsAPI{}
	toTest := c.checkValidRequest(c.handleAvgClicks(api))
	toTest.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("HandleAvgMetrics.ServeHTTP result status code: expected %d, got %d", http.StatusOK, w.Result().StatusCode)
	}
}

func TestHandleAvgMetricsNoAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "/groups/Bk1hmwBHfQK/countries/averages", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	c := BitlyClientInfo{}
	api := &bitlinksMetricsAPI{}
	toTest := c.checkValidRequest(c.handleAvgClicks(api))
	toTest.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusForbidden {
		t.Errorf("HandleAvgMetricsNoAuth.ServeHTTP result status code: expected %d, got %d", http.StatusForbidden, w.Result().StatusCode)
	}

}

func TestHandleAvgMetricsInternalError(t *testing.T) {
	req, err := http.NewRequest("GET", "/groups/Bk1hmwBHfQK/countries/averages", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	c := BitlyClientInfo{}
	mock := &mockBitlinksMetricsAPI{}
	mock.getBitlinksForGroupError = errors.New("Could not retrieve bitlinks for the requested group")
	toTest := c.handleAvgClicks(mock)
	toTest.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("HandleAvgMetricsInternalError.ServeHTTP result status code: expected %d, got %d", http.StatusInternalServerError, w.Result().StatusCode)
	}

}

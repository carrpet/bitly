package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockBitlinksMetricsAPI struct{}

func (m *mockBitlinksMetricsAPI) GetUserInfo(c bitlyClient) (*bitlyUserInfo, error) {
	return &bitlyUserInfo{GroupGuid: "abcdefgh", Name: "petertest"}, nil

}

func (m *mockBitlinksMetricsAPI) GetBitlinksForGroup(c bitlyClient, guid string) (*bitlyGroupsBitLinks, error) {
	l := []Bitlink{{Link: "http://something", ID: "something"},
		{Link: "http://something1", ID: "something1"},
		{Link: "http://something2", ID: "something2"},
	}
	return &bitlyGroupsBitLinks{Links: l}, nil

}

func (m *mockBitlinksMetricsAPI) GetClicksByCountry(c bitlyClient, b Bitlink) (*ClickMetrics, error) {
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

	expected := []struct {
		country        string
		expectedClicks int
	}{
		{"US", 80 * 3},
		{"Argentina", 40 * 3},
		{"Norway", 10 * 3},
		{"Sweden", 60 * 3},
		{"China", 1200 * 3},
	}
	context := BitlyClientInfo{}
	mock := &mockBitlinksMetricsAPI{}

	cc, err := context.avgClicks(mock)
	if err != nil {
		t.Error("We fucked up!")
	}
	if len(cc) != 5 {
		t.Errorf("avgClicks return value has length %d, expected length 5", len(cc))
	}
	for i, item := range cc {
		if expected[i].country != item.Country || expected[i].expectedClicks != item.Clicks {
			t.Errorf("avgClicks return value: expected country %s and clicks %d, received country %s and clicks %d",
				expected[i].country, expected[i].expectedClicks, item.Country, item.Clicks)
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
	api := &BitlinksMetricsAPI{}
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
	api := &BitlinksMetricsAPI{}
	toTest := c.checkValidRequest(c.handleAvgClicks(api))
	toTest.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusForbidden {
		t.Errorf("HandleAvgMetricsNoAuth.ServeHTTP result status code: expected %d, got %d", http.StatusForbidden, w.Result().StatusCode)
	}

}

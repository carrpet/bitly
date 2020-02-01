package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockBitlinksMetricsAPI struct{}

func (m *mockBitlinksMetricsAPI) GetUserInfo(c bitlyClient) (*bitlyUserInfo, error) {
	return nil, nil

}

func (m *mockBitlinksMetricsAPI) GetBitlinksForGroup(c bitlyClient, guid string) (*bitlyGroupsBitLinks, error) {
	return nil, nil

}

func (m *mockBitlinksMetricsAPI) GetClicksByCountry(c bitlyClient, b Bitlink) (*ClickMetrics, error) {
	return nil, nil

}

// Main handler tests

/*
func TestAvgClickMetricsHandler(t *testing.T) {
	context := BitlyClientInfo{}
	mock := &mockBitlinksMetricsAPI{}

	_, err := context.avgClicks(mock)
	if err != nil {
		t.Error("We fucked up!")
	}

}
*/

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

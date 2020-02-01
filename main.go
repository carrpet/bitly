package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// The problem:

// To get all the bitlinks for a user's default group then for each bitlink
// get the number of user clicks by country, add all the clicks up by country
// then we divide those numbers by 30 and return a list of countries and average
//clicks

//Pseudocode and naive implementation:
/*
countrytoclickshashtable(ctcht) = {}
groupId = GET(user)
bitlinks = GET(groups/{groupId}/bitlinks)
for each link in bitlinks:
     clicksbycountry = GET(bitlinks/{link}/countries)
     ctcht{clicksbycountry[country]} += clicksbycountry[clicks]
return map(divideclicksby30,ctcht)

*/

// use this for testing
//token := "5ad8274a49bcd964f23d4b685c272c37de718711"

func (c *BitlyClientInfo) checkValidRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No authorization header provided", http.StatusForbidden)
			return
		}
		auths := strings.Split(authHeader, " ")
		if len(auths) < 2 {
			http.Error(w, "No bearer token provided", http.StatusForbidden)
			return
		}
		c.Token = auths[1]
		h(w, r)

	}
}
func (c *BitlyClientInfo) handleAvgClicks(api BitlinksMetrics) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		type avgClicksResponse struct {
			Units   int            `json:"units"`
			Facet   string         `json:"facet"`
			UnitRef time.Time      `json:"unit_reference"`
			Unit    int            `json:"unit"`
			Metrics []CountryClick `json:"metrics"`
		}

		res, err := c.avgClicks(api)

		if err != nil {
			panic(err)
			//TODO: return proper http error code
		}

		data := avgClicksResponse{
			Facet:   "countries",
			UnitRef: time.Now(),
			Metrics: res,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)

	}

}

func (c *BitlyClientInfo) avgClicks(api BitlinksMetrics) ([]CountryClick, error) {
	clicksByCountry := map[string]int{}

	userInfo, err := api.GetUserInfo(c)
	if err != nil {
		panic(err)
	}

	// retrieve all the links and stick into a hashtable
	grouplinks, err := api.GetBitlinksForGroup(c, userInfo.GroupGuid)
	if err != nil {
		panic(err)
	}

	var cc *ClickMetrics
	for i := 0; i < len(grouplinks.Links); i++ {
		cc, err := api.GetClicksByCountry(c, grouplinks.Links[i])
		if err != nil {
			panic(err)
		}
		for j := 0; j < len(cc.Metrics); j++ {
			_, ok := clicksByCountry[cc.Metrics[j].Country]
			if !ok {
				clicksByCountry[cc.Metrics[j].Country] = cc.Metrics[j].Clicks
			} else {
				clicksByCountry[cc.Metrics[j].Country] += cc.Metrics[j].Clicks
			}
			fmt.Printf("Clicks By Country: clicks: %d, country: %s\n", cc.Metrics[j].Clicks, cc.Metrics[j].Country)
		}
	}

	if cc != nil {
		return computeAvgClicks(&cc.Metrics), nil
	} else {
		return []CountryClick{}, nil
	}

}

func computeAvgClicks(cc *[]CountryClick) []CountryClick {
	for _, val := range *cc {
		val.Clicks = val.Clicks / 30
	}
	return *cc
}

func main() {
	context := BitlyClientInfo{}
	api := &BitlinksMetricsAPI{}
	http.HandleFunc("/groups/{groupGuid}/countries/averages", context.checkValidRequest(context.handleAvgClicks(api)))
	log.Fatal(http.ListenAndServe(":8080", nil))

}

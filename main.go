package main

import (
	"bitly/client"
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
func handleAvgClicks() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		type avgClicksResponse struct {
			Units   int                   `json:"units"`
			Facet   string                `json:"facet"`
			UnitRef time.Time             `json:"unit_reference"`
			Unit    int                   `json:"unit"`
			Metrics []client.CountryClick `json:"metrics"`
		}
		authHeader := r.Header.Get("Authorization")
		authstrings := strings.Split(authHeader, " ")
		token := authstrings[1]
		bc := &client.BitlyClientInfo{Token: token}
		clicksByCountry := map[string]int{}

		userInfo, err := client.GetUserInfo(bc)
		if err != nil {
			panic(err)
		}

		// retrieve all the links and stick into a hashtable
		grouplinks, err := client.GetBitlinksForGroup(bc, userInfo.GroupGuid)
		if err != nil {
			panic(err)
		}

		var cc *client.ClickMetrics
		for i := 0; i < len(grouplinks.Links); i++ {
			cc, err := client.GetClicksByCountry(bc, grouplinks.Links[i])
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

		var metrics []client.CountryClick
		if cc != nil {
			metrics = avgClicks(&cc.Metrics)
		} else {
			metrics = []client.CountryClick{}
		}
		data := avgClicksResponse{
			Facet:   "countries",
			UnitRef: time.Now(),
			Metrics: metrics,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)

	}

}

func avgClicks(cc *[]client.CountryClick) []client.CountryClick {
	for _, val := range *cc {
		val.Clicks = val.Clicks / 30
	}
	return *cc
}

func main() {

	http.HandleFunc("/groups/{groupGuid}/countries/averages", handleAvgClicks())
	log.Fatal(http.ListenAndServe(":8080", nil))

}

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
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

const HTTP_TIMEOUT = 10 * time.Second

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

// avgClicks does the brunt of the handler's work, calling the API methods,
// aggregating the results and computing the average values for the response.
func (c *BitlyClientInfo) avgClicks(api BitlinksMetrics) ([]CountryClick, error) {
	clicksByCountry := map[string]int{}

	userInfo, err := api.GetUserInfo(c)
	if err != nil {
		return nil, err
	}

	// retrieve all the links and stick into a hashtable
	grouplinks, err := api.GetBitlinksForGroup(c, userInfo.GroupGuid)
	if err != nil {
		return nil, err
	}

	for _, link := range grouplinks.Links {
		cc, err := api.GetBitlinkClicksByCountry(c, link)
		if err != nil {
			return nil, err
		}
		for _, m := range cc.Metrics {
			_, ok := clicksByCountry[m.Country]
			if !ok {
				clicksByCountry[m.Country] = m.Clicks
			} else {
				clicksByCountry[m.Country] += m.Clicks
			}
		}
	}

	if len(clicksByCountry) > 0 {
		arr := toCountryClickArray(clicksByCountry)
		return computeAvgClicks(&arr, DEFAULT_DAYS), nil
	} else {
		return []CountryClick{}, nil
	}

}

func computeAvgClicks(cc *[]CountryClick, days int) []CountryClick {
	for _, val := range *cc {
		val.Clicks = val.Clicks / days
	}
	return *cc
}

func toCountryClickArray(cc map[string]int) []CountryClick {
	ret := []CountryClick{}
	for k, v := range cc {
		ret = append(ret, CountryClick{Clicks: v, Country: k})
	}
	return ret
}

func handleTest() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello there!`))
	}
}

func main() {
	clientInfo := BitlyClientInfo{}
	api := &bitlinksMetricsAPI{}
	r := mux.NewRouter()
	r.HandleFunc("/groups/{groupGuid}/countries/averages", clientInfo.checkValidRequest(clientInfo.handleAvgClicks(api)))
	r.HandleFunc("/test", handleTest())
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  HTTP_TIMEOUT,
		WriteTimeout: HTTP_TIMEOUT,
	}

	//goroutine to start the server
	go func() {
		log.Println("Starting Bitly Metrics server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Printf("Bitly Metrics server listening on: %s", srv.Addr)
	awaitShutdown(srv)
}

func awaitShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
	ctx, cancel := context.WithTimeout(context.Background(), HTTP_TIMEOUT)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down Bitly Metrics server")
	os.Exit(0)
}

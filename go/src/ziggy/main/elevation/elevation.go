package main

import (
	"flag"
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"os"
)

var t = flag.String("t", "", "api-token")

type results struct {
	Results []struct {
		Elevation float64 `json:"elevation"`
		Location  struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
		Resolution float64 `json:"resolution"`
	} `json:"results"`
	Status string `json:"status"`
}

func main() {
	flag.Parse()
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/elevation/json?locations=%s,%s&key=%s", flag.Arg(0), flag.Arg(1), *t)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	var data results
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%v\n", data.Results[0].Elevation)

}

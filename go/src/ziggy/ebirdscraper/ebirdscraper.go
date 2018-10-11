package ebirdscraper

import (
	"net/http"
	"fmt"
	"encoding/json"
)

func FetchObsOnDayInRegion(y, m, d int, region, apiKey string) ([]*Observation, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://ebird.org/ws2.0/data/obs/%s/historic/%d/%d/%d?rank=mrec&detail=full&cat=species",region, y, m, d),
		nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("X-eBirdApiToken", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}


	var out []*Observation
	return out, json.NewDecoder(resp.Body).Decode(&out)
}


func FilterForSpecies(obs []*Observation, species string) []*Observation {
	var out []*Observation
	for i, v := range obs {
		if v.SpeciesCode == species {
			out = append(out, obs[i])
		}
	}
	
	return out
}

type Observation struct {
	SpeciesCode      string  `json:"speciesCode"`
	ComName          string  `json:"comName"`
	SciName          string  `json:"sciName"`
	LocID            string  `json:"locId"`
	LocName          string  `json:"locName"`
	ObsDt            string  `json:"obsDt"`
	HowMany          int     `json:"howMany"`
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
	ObsValid         bool    `json:"obsValid"`
	ObsReviewed      bool    `json:"obsReviewed"`
	LocationPrivate  bool    `json:"locationPrivate"`
	Subnational2Code string  `json:"subnational2Code"`
	Subnational2Name string  `json:"subnational2Name"`
	Subnational1Code string  `json:"subnational1Code"`
	Subnational1Name string  `json:"subnational1Name"`
	CountryCode      string  `json:"countryCode"`
	CountryName      string  `json:"countryName"`
	UserDisplayName  string  `json:"userDisplayName"`
	SubID            string  `json:"subId"`
	ObsID            string  `json:"obsId"`
	ChecklistID      string  `json:"checklistId"`
	PresenceNoted    bool    `json:"presenceNoted"`
	HasComments      bool    `json:"hasComments"`
	HasRichMedia     bool    `json:"hasRichMedia"`
	LastName         string  `json:"lastName"`
	FirstName        string  `json:"firstName"`
}
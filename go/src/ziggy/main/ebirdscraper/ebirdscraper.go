package main

import (
	"flag"
	"time"
	"encoding/csv"
	"os"
	"reflect"
	"ziggy/ebirdscraper"
	"log"
	"fmt"
	"sync"
)

const flagTimeFmt = "1/2/2006"

var apiKey = flag.String("t", "", "api token")
var startDate = flag.String("s", "", "start date m/d/yyyy")
var endDate = flag.String("e", "", "end date mm/d/yyyy")
var region = flag.String("r", "US-WA", "region")
var species = flag.String("sp", "soogro1", "species")

func main() {
	flag.Parse()

	start, err := time.Parse(flagTimeFmt, *startDate)
	if err != nil {
		log.Fatal(err)
	}

	end, err := time.Parse(flagTimeFmt, *endDate)
	if err != nil {
		log.Fatal(err)
	}


	w := csv.NewWriter(os.Stdout)

	obsType := reflect.TypeOf(ebirdscraper.Observation{})

	var headers []string
	for i := 0; i < obsType.NumField(); i++ {
		headers = append(headers, obsType.Field(i).Name)
	}

	w.Write(headers)
	w.Flush()

	lock := &sync.Mutex{}
	write := func (rec []string) {
		lock.Lock()
		w.Write(rec)
		w.Flush()
		lock.Unlock()
	}

	worker := func(q chan time.Time) {
		for next := range q {
			obsType := reflect.TypeOf(ebirdscraper.Observation{})
			observations, err := ebirdscraper.FetchObsOnDayInRegion(next.Year(), int(next.Month()), next.Day(), *region, *apiKey)
			if err != nil {
				log.Fatal(err)
			}

			specObs := ebirdscraper.FilterForSpecies(observations, *species)

			for _, so := range specObs {
				soV := reflect.ValueOf(so).Elem()
				var row []string
				for i := 0; i < obsType.NumField(); i++ {
					row = append(row, fmt.Sprintf("%v", soV.Field(i).Interface()))
				}

				write(row)
			}
		}
	}


	next := start
	q := make(chan time.Time, 1000)
	for i := 0; i < 20; i++ {
		go worker(q)
	}

	for next.Before(end) {
		q <- next
		next = next.Add(24*time.Hour)
	}
	close(q)

}

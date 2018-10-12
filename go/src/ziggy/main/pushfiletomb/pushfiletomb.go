package main

import (
	"flag"
	"log"
	"os"
	"ziggy/mapbox"
)

var token = flag.String("t", os.Getenv("MAPBOX_API_TOKEN"), "mapbox api token")
var user = flag.String("u", "hdh", "mapbox username")

func main() {
	flag.Parse()
	client := mapbox.NewClient(*token, *user)
	if err := client.CreateOrReplaceTileset(flag.Arg(1), flag.Arg(0)); err != nil {
		log.Fatal(err)
	}
}

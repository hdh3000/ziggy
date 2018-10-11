package main

import (
	"flag"
	"log"
	"path"
	"strings"
	"time"
	"ziggy/exporter"
	"ziggy/mapbox"
)

var mbUser = flag.String("u", "hdh", "mapbox username")
var mbKey = flag.String("k", "", "mapbox api key")
var mdStore = flag.String("md", "/Users/hdh/src/ziggy/etc/sqlexports.json", "where to log data exports")

func main() {
	flag.Parse()
	file := flag.Arg(0)
	query := flag.Arg(1)
	exportPath, err := exporter.ExportShp(query, file)
	if err != nil {
		log.Fatal(err)
	}

	mb := mapbox.NewClient(*mbKey, *mbUser)

	tilesetName := strings.Split(path.Base(file), ".")[0]

	if len(tilesetName) > 32 {
		tilesetName = tilesetName[:32] // that is max length
	}

	if err := mb.CreateOrReplaceTileset(exportPath, tilesetName); err != nil {
		log.Fatal(err)
	}

	if err := exporter.WriteExportMetadata(*mdStore, &exporter.SQLExportMetaData{
		Date:  time.Now(),
		Query: query,
		Name:  tilesetName,
	}); err != nil {
		log.Fatal(err)
	}

}

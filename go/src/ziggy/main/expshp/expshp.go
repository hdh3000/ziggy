package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
	"ziggy/exporter"
	"ziggy/filedata"
	"ziggy/mapbox"
)

var mbUser = flag.String("u", "hdh", "mapbox username")
var mbKey = flag.String("t", os.Getenv("MAPBOX_API_TOKEN"), "mapbox api key")
var mdStore = flag.String("md", "/Users/hdh/src/ziggy/etc/sqlexports.json", "where to log data exports")

func main() {
	flag.Parse()
	file := flag.Arg(0)

	queryFile := flag.Arg(1)

	qf, err := os.Open(queryFile)
	if err != nil {
		log.Fatal(err)
	}
	defer qf.Close()

	query, _ := ioutil.ReadAll(qf)

	exportPath, err := exporter.ExportShp(string(query), file)
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

	if err := filedata.NewMgr(*mdStore).Put(&exporter.SQLExportMetaData{
		Date:  time.Now(),
		Query: string(query),
		Name:  tilesetName,
	}); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"flag"
	"log"
	"ziggy/filedata"
	"ziggy/importer"
)

var mdStore = flag.String("md", "/Users/hdh/src/ziggy/etc/sqlimports.json", "what to use for the metadata store?")
var fType = flag.String("ft", "shp", "file type, can be rast | shp")

func main() {
	flag.Parse()
	systemName := flag.Arg(0)
	filePath := flag.Arg(1)
	url := flag.Arg(2)

	var ds importer.SQLDataSource
	switch *fType {
	case "shp":
		ds = importer.NewShpFile(systemName, url, filePath)
	case "rast":
		ds = importer.NewRaster(systemName, url, filePath)
	default:
		log.Fatalf("%q is not a recognized file type", *fType)
	}

	metadata, err := ds.ImpToSql()
	if err != nil {
		log.Fatal(err)
	}

	if err := filedata.NewMgr(*mdStore).Put(metadata); err != nil {
		log.Fatal(err)
	}
}

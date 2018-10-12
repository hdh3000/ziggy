package main

import (
	"flag"
	"log"
	"ziggy/filedata"
	"ziggy/importer"
)

var mdStore = flag.String("md", "/Users/hdh/src/ziggy/etc/sqlimports.json", "what to use for the metadata store?")

func main() {
	flag.Parse()
	systemName := flag.Arg(0)
	filePath := flag.Arg(1)
	url := flag.Arg(2)

	shp := importer.NewShpFile(systemName, url, filePath)

	metadata, err := shp.ImpToSql()
	if err != nil {
		log.Fatal(err)
	}

	if err := filedata.NewMgr(*mdStore).Put(metadata); err != nil {
		log.Fatal(err)
	}
}

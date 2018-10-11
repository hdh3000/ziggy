package main

import (
	"flag"
	"log"
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

	if err := importer.WriteSqlMetaData(*mdStore, metadata); err != nil {
		log.Fatal(err)
	}

}

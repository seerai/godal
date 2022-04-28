package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	gdal "github.com/seerai/godal"
)

func main() {
	flag.Parse()
	inputFile := flag.Arg(0)
	fmt.Printf("Input filename: %s\n", inputFile)

	b, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	vf, err := gdal.VSIFileFromMemBuffer("/vsimem/test.tif", b, false)
	if err != nil {
		panic(err)
	}
	ds, err := gdal.Open("/vsimem/test.tif", gdal.ReadOnly)
	if err != nil {
		log.Fatal(err)
	}
	defer ds.Close()

	fmt.Println(ds.GeoTransform())
	fmt.Println(ds.Driver().LongName())

	defer gdal.VSIFCloseL(vf)
	defer gdal.VSIUnlink("/vsimem/test.tif")
	
}

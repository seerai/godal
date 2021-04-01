package main

import (
	"fmt"

	"github.com/seerai/godal"
)

func main() {

	//var ds gdal.Dataset
	imageList := []string{}
	options := []string{}

	outputFile := ""

	outputDs := gdal.BuildVRT(outputFile, imageList, options)

	fmt.Println(outputDs)

	outputDs.Close()

	//outputDs := gdal.GDALWarp(outputFile, gdal.Dataset{}, []gdal.Dataset{ds}, options)
	fmt.Printf("End program\n")
}

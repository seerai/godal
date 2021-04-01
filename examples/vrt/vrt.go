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

	fmt.Printf("End program\n")
}

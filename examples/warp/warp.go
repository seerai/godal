package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/seerai/godal"
)

func main() {
	flag.Parse()
	inputFile := flag.Arg(0)
	outputFile := flag.Arg(1)
	if len(flag.Args()) < 3 {
		fmt.Printf("Usage: test_warp inputFile outputFile options\n")
		return
	}
	options := flag.Args()[2:]
	if inputFile == "" {
		fmt.Printf("Usage: test_warp inputFile outputFile options\n")
		return
	}
	fmt.Printf("Input filename: %s\n", inputFile)
	if outputFile == "" {
		fmt.Printf("Usage: test_warp inputFile outputFile options\n")
		return
	}
	fmt.Printf("Output filename: %s\n", outputFile)

	fmt.Printf("Warp options: %s\n", strings.Join(options, " "))

	ds, err := gdal.Open(inputFile, gdal.ReadOnly)
	if err != nil {
		log.Fatal(err)
	}
	defer ds.Close()

	outputDs, err := gdal.Warp(outputFile, gdal.Dataset{}, []gdal.Dataset{ds}, options)
	defer outputDs.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("End program\n")
}

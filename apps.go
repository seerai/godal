package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"
#include "gdal_utils.h"

#cgo linux  pkg-config: gdal
#cgo darwin pkg-config: gdal
#cgo windows LDFLAGS: -Lc:/gdal/release-1600-x64/lib -lgdal_i
#cgo windows CFLAGS: -IC:/gdal/release-1600-x64/include
*/
import "C"
import (
	"fmt"
	"unsafe"
)

var _ = fmt.Println

/* --------------------------------------------- */
/* GDAL utilities                                */
/* --------------------------------------------- */

// TranslateOptions holds options to be passed to gdal translated
type TranslateOptions struct {
	cval *C.GDALTranslateOptions
}

// WarpAppOptions holds options to be passed to gdal translated
type WarpAppOptions struct {
	cval *C.GDALWarpAppOptions
}

//RasterizeAppOptions options options to be passed to gdal
type RasterizeAppOptions struct {
	cval *C.GDALRasterizeOptions
}

// Translate is a utility to convert images into different formats
func Translate(
	destName string,
	srcDS Dataset,
	options []string,
) Dataset {

	var err C.int

	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	gdalTranslateOptions := TranslateOptions{C.GDALTranslateOptionsNew((**C.char)(unsafe.Pointer(&cOptions[0])), nil)}

	outputDs := C.GDALTranslate(
		C.CString(destName),
		srcDS.cval,
		gdalTranslateOptions.cval,
		&err,
	)

	return Dataset{outputDs}

}

// Warp is a utility to warp images into different projections
func Warp(
	destName string,
	dstDs Dataset,
	srcDs []Dataset,
	options []string,
) (Dataset, error) {

	var err C.int

	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	gdalWarpOptions := WarpAppOptions{C.GDALWarpAppOptionsNew((**C.char)(unsafe.Pointer(&cOptions[0])), nil)}
	if gdalWarpOptions.cval == nil {
		fmt.Println("GDALWarpAppOptionsNew() returned a null pointer.")
		return Dataset{}, ErrFailure
	}

	pahSrcDs := make([]C.GDALDatasetH, len(srcDs)+1)
	for i := 0; i < len(srcDs); i++ {
		pahSrcDs[i] = srcDs[i].cval
	}
	pahSrcDs[len(srcDs)] = (C.GDALDatasetH)(unsafe.Pointer(nil))

	outputDs := C.GDALWarp(
		C.CString(destName),
		dstDs.cval,
		C.int(len(srcDs)),
		(*C.GDALDatasetH)(unsafe.Pointer(&pahSrcDs[0])),
		gdalWarpOptions.cval,
		&err,
	)

	if err != 0 {
		return Dataset{outputDs}, ErrFailure
	}

	return Dataset{outputDs}, nil

}

// BuildVRT creates a new dataset that is the mosaic of the input files.
func BuildVRT(
	outputFile string,
	inputDatasets []string,
	options []string,
) Dataset {

	// Flag to store error code
	var err C.int

	// Parse the user options
	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	buildVRTOptions := BuildVRTOptions{C.GDALBuildVRTOptionsNew((**C.char)(unsafe.Pointer(&cOptions[0])), nil)}

	// Output file path for the VRT dataset
	cPath := C.CString(outputFile)
	defer C.free(unsafe.Pointer(cPath))

	// Create C strings for the input files
	length = len(inputDatasets)
	srcDSNames := make([]*C.char, length+1)

	for i := 0; i < length; i++ {
		srcDSNames[i] = C.CString(inputDatasets[i])
		defer C.free(unsafe.Pointer(srcDSNames[i]))
	}
	srcDSNames[length] = (*C.char)(unsafe.Pointer(nil))

	// Call the BuildVRT function
	outputDs := C.GDALBuildVRT(
		cPath,                     // Output dataset path
		C.int(len(inputDatasets)), // Number of input datasets
		nil,                       // pointer to input dataset (nil)
		(**C.char)(unsafe.Pointer(&srcDSNames[0])),
		buildVRTOptions.cval,
		&err,
	)

	return Dataset{outputDs}
}

// Rasterize creates a new dataset that is the rasterization of input features.
func Rasterize(outputDest string, outputDataset Dataset, inputDataset Dataset, options []string) (Dataset, error) {
	var err C.int

	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	rasterizeOptions := RasterizeAppOptions{C.GDALRasterizeOptionsNew((**C.char)(unsafe.Pointer(&cOptions[0])), nil)}
	if rasterizeOptions.cval == nil {
		fmt.Println("GDALRasterizeOptionsNew() returned a null pointer.")
		return Dataset{}, ErrFailure
	}
	defer C.GDALRasterizeOptionsFree(rasterizeOptions.cval)

	var outputDs C.GDALDatasetH
	if outputDest != "" {
		outputDs = C.GDALRasterize(
			C.CString(outputDest),
			outputDataset.cval,
			inputDataset.cval,
			rasterizeOptions.cval,
			&err,
		)
	} else {
		outputDs = C.GDALRasterize(
			nil,
			outputDataset.cval,
			inputDataset.cval,
			rasterizeOptions.cval,
			&err,
		)
	}

	if err != 0 {
		return Dataset{outputDs}, ErrFailure
	}

	return Dataset{outputDs}, nil

}

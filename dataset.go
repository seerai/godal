package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"

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

// Driver gets the driver to which this dataset relates
func (dataset Dataset) Driver() Driver {
	driver := Driver{C.GDALGetDatasetDriver(dataset.cval)}
	return driver
}

// FileList fetches files forming the dataset.
func (dataset Dataset) FileList() []string {
	p := C.GDALGetFileList(dataset.cval)
	var strings []string
	q := uintptr(unsafe.Pointer(p))
	for {
		p = (**C.char)(unsafe.Pointer(q))
		if *p == nil {
			break
		}
		strings = append(strings, C.GoString(*p))
		q += unsafe.Sizeof(q)
	}

	return strings
}

// Close closes the dataset
func (dataset Dataset) Close() {
	C.GDALClose(dataset.cval)
	return
}

// RasterXSize fetches X size of raster
func (dataset Dataset) RasterXSize() int {
	xSize := int(C.GDALGetRasterXSize(dataset.cval))
	return xSize
}

// RasterYSize fetches Y size of raster
func (dataset Dataset) RasterYSize() int {
	ySize := int(C.GDALGetRasterYSize(dataset.cval))
	return ySize
}

// RasterCount fetches the number of raster bands in the dataset
func (dataset Dataset) RasterCount() int {
	count := int(C.GDALGetRasterCount(dataset.cval))
	return count
}

// RasterBand fetches a raster band object from a dataset
func (dataset Dataset) RasterBand(band int) RasterBand {
	rasterBand := RasterBand{C.GDALGetRasterBand(dataset.cval, C.int(band))}
	return rasterBand
}

// AddBand adds a band to a dataset
func (dataset Dataset) AddBand(dataType DataType, options []string) error {
	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALAddBand(
		dataset.cval,
		C.GDALDataType(dataType),
		(**C.char)(unsafe.Pointer(&cOptions[0])),
	).Err()
}

// AutoCreateWarpedVRT creates a warped VRT from a source to destination projection specified as WKT
func (dataset Dataset) AutoCreateWarpedVRT(srcWKT, dstWKT string, resampleAlg ResampleAlg) (Dataset, error) {
	c_srcWKT := C.CString(srcWKT)
	defer C.free(unsafe.Pointer(c_srcWKT))
	c_dstWKT := C.CString(dstWKT)
	defer C.free(unsafe.Pointer(c_dstWKT))

	h := C.GDALAutoCreateWarpedVRT(dataset.cval, c_srcWKT, c_dstWKT, C.GDALResampleAlg(resampleAlg), 0.0, nil)
	d := Dataset{h}
	if h == nil {
		return d, fmt.Errorf("AutoCreateWarpedVRT failed")
	}
	return d, nil

}

// Unimplemented: GDALBeginAsyncReader
// Unimplemented: GDALEndAsyncReader

// IO reads / writes a region of image data from multiple bands
func (dataset Dataset) IO(
	rwFlag RWFlag,
	xOff, yOff, xSize, ySize int,
	buffer interface{},
	bufXSize, bufYSize int,
	bandCount int,
	bandMap []int,
	pixelSpace, lineSpace, bandSpace int,
	readRawBytes bool,
) error {
	var dataType DataType
	var dataPtr unsafe.Pointer
	if readRawBytes {
		data := ([]uint8)(buffer.([]uint8))
		dataType = dataset.RasterBand(1).RasterDataType()
		dataPtr = unsafe.Pointer(&data[0])
	} else {
		switch data := buffer.(type) {
		case []int8:
			dataType = Byte
			dataPtr = unsafe.Pointer(&data[0])
		case []uint8:
			dataType = Byte
			dataPtr = unsafe.Pointer(&data[0])
		case []int16:
			dataType = Int16
			dataPtr = unsafe.Pointer(&data[0])
		case []uint16:
			dataType = UInt16
			dataPtr = unsafe.Pointer(&data[0])
		case []int32:
			dataType = Int32
			dataPtr = unsafe.Pointer(&data[0])
		case []uint32:
			dataType = UInt32
			dataPtr = unsafe.Pointer(&data[0])
		case []float32:
			dataType = Float32
			dataPtr = unsafe.Pointer(&data[0])
		case []float64:
			dataType = Float64
			dataPtr = unsafe.Pointer(&data[0])
		default:
			return fmt.Errorf("Error: buffer is not a valid data type (must be a valid numeric slice)")
		}
	}

	return C.GDALDatasetRasterIO(
		dataset.cval,
		C.GDALRWFlag(rwFlag),
		C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize),
		dataPtr,
		C.int(bufXSize), C.int(bufYSize),
		C.GDALDataType(dataType),
		C.int(bandCount),
		(*C.int)(unsafe.Pointer(&IntSliceToCInt(bandMap)[0])),
		C.int(pixelSpace), C.int(lineSpace), C.int(bandSpace),
	).Err()
}

// BasicRead reads from a dataset with some basic and very common assumptions
func (dataset Dataset) BasicRead(
	xOff, yOff, xSize, ySize int,
	bands []int,
	buffer interface{},
) error {

	bufXSize := xSize
	bufYSize := ySize
	pixelSpace := 0
	lineSpace := 0
	bandSpace := 0
	readRawBytes := true

	return dataset.IO(
		Read,
		xOff, yOff, xSize, ySize,
		buffer,
		bufXSize, bufYSize,
		len(bands),
		bands,
		pixelSpace, lineSpace, bandSpace,
		readRawBytes,
	)
}

// BasicWrite writes to a dataset with some basic and very common assumptions
func (dataset Dataset) BasicWrite(
	xOff, yOff, xSize, ySize int,
	bands []int,
	buffer interface{},
) error {

	bufXSize := xSize
	bufYSize := ySize
	pixelSpace := 0
	lineSpace := 0
	bandSpace := 0
	readRawBytes := true

	return dataset.IO(
		Write,
		xOff, yOff, xSize, ySize,
		buffer,
		bufXSize, bufYSize,
		len(bands),
		bands,
		pixelSpace, lineSpace, bandSpace,
		readRawBytes,
	)
}

// AdviseRead advises driver of upcoming read requests
func (dataset Dataset) AdviseRead(
	rwFlag RWFlag,
	xOff, yOff, xSize, ySize, bufXSize, bufYSize int,
	dataType DataType,
	bandCount int,
	bandMap []int,
	options []string,
) error {
	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALDatasetAdviseRead(
		dataset.cval,
		C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize),
		C.int(bufXSize), C.int(bufYSize),
		C.GDALDataType(dataType),
		C.int(bandCount),
		(*C.int)(unsafe.Pointer(&IntSliceToCInt(bandMap)[0])),
		(**C.char)(unsafe.Pointer(&cOptions[0])),
	).Err()
}

// Projection fetches the projection definition string for this dataset
func (dataset Dataset) Projection() string {
	proj := C.GoString(C.GDALGetProjectionRef(dataset.cval))
	return proj
}

// SetProjection sets the projection reference string
func (dataset Dataset) SetProjection(proj string) error {
	cProj := C.CString(proj)
	defer C.free(unsafe.Pointer(cProj))

	return C.GDALSetProjection(dataset.cval, cProj).Err()
}

// GeoTransform gets the affine transformation coefficients
func (dataset Dataset) GeoTransform() [6]float64 {
	var transform [6]float64
	C.GDALGetGeoTransform(dataset.cval, (*C.double)(unsafe.Pointer(&transform[0])))
	return transform
}

// SetGeoTransform sets the affine transformation coefficients
func (dataset Dataset) SetGeoTransform(transform [6]float64) error {
	return C.GDALSetGeoTransform(
		dataset.cval,
		(*C.double)(unsafe.Pointer(&transform[0])),
	).Err()
}

// InvGeoTransform returns the inverted transform
func (dataset Dataset) InvGeoTransform() [6]float64 {
	return InvGeoTransform(dataset.GeoTransform())
}

// InvGeoTransform inverts the supplied transform
func InvGeoTransform(transform [6]float64) [6]float64 {
	var result [6]float64
	C.GDALInvGeoTransform((*C.double)(unsafe.Pointer(&transform[0])), (*C.double)(unsafe.Pointer(&result[0])))
	return result
}

// GetGCPCount gets number of GCPs
func (dataset Dataset) GetGCPCount() int {
	count := C.GDALGetGCPCount(dataset.cval)
	return int(count)
}

// Unimplemented: GDALGetGCPProjection
// Unimplemented: GDALGetGCPs
// Unimplemented: GDALSetGCPs

// GetInternalHandle fetches a format specific internally meaningful handle
func (dataset Dataset) GetInternalHandle(request string) unsafe.Pointer {
	cRequest := C.CString(request)
	defer C.free(unsafe.Pointer(cRequest))

	ptr := C.GDALGetInternalHandle(dataset.cval, cRequest)
	return ptr
}

// ReferenceDataset adds one to dataset reference count
func (dataset Dataset) ReferenceDataset() int {
	count := C.GDALReferenceDataset(dataset.cval)
	return int(count)
}

// DereferenceDataset subtracts one from dataset reference count
func (dataset Dataset) DereferenceDataset() int {
	count := C.GDALDereferenceDataset(dataset.cval)
	return int(count)
}

// BuildOverviews builds raster overview(s)
func (dataset Dataset) BuildOverviews(
	resampling string,
	nOverviews int,
	overviewList []int,
	nBands int,
	bandList []int,
	progress ProgressFunc,
	data interface{},
) error {
	cResampling := C.CString(resampling)
	defer C.free(unsafe.Pointer(cResampling))

	arg := &goGDALProgressFuncProxyArgs{progress, data}

	return C.GDALBuildOverviews(
		dataset.cval,
		cResampling,
		C.int(nOverviews),
		(*C.int)(unsafe.Pointer(&IntSliceToCInt(overviewList)[0])),
		C.int(nBands),
		(*C.int)(unsafe.Pointer(&IntSliceToCInt(bandList)[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Unimplemented: GDALGetOpenDatasets

// Access returns access flag
func (dataset Dataset) Access() Access {
	accessVal := C.GDALGetAccess(dataset.cval)
	return Access(accessVal)
}

// FlushCache writes all write cached data to disk
func (dataset Dataset) FlushCache() {
	C.GDALFlushCache(dataset.cval)
	return
}

// CreateMaskBand adds a mask band to the dataset
func (dataset Dataset) CreateMaskBand(flags int) error {
	return C.GDALCreateDatasetMaskBand(dataset.cval, C.int(flags)).Err()
}

// CopyWholeRaster copies all dataset raster data
func (sourceDataset Dataset) CopyWholeRaster(
	destDataset Dataset,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{progress, data}

	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALDatasetCopyWholeRaster(
		sourceDataset.cval,
		destDataset.cval,
		(**C.char)(unsafe.Pointer(&cOptions[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

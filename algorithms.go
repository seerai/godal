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
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

var _ = fmt.Println

const (
	RPC_LINE_NUM_COEFF = "LINE_NUM_COEFF"
	RPC_LINE_DEN_COEFF = "LINE_DEN_COEFF"
	RPC_SAMP_NUM_COEFF = "SAMP_NUM_COEFF"
	RPC_SAMP_DEN_COEFF = "SAMP_DEN_COEFF"
)

/* --------------------------------------------- */
/* Misc functions                                */
/* --------------------------------------------- */

// ComputeMedianCutPCT computes optimal PCT for RGB image
func ComputeMedianCutPCT(
	red, green, blue RasterBand,
	colors int,
	ct ColorTable,
	progress ProgressFunc,
	data interface{},
) int {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	err := C.GDALComputeMedianCutPCT(
		red.cval,
		green.cval,
		blue.cval,
		nil,
		C.int(colors),
		ct.cval,
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	)
	return int(err)
}

// 24bit to 8bit conversion with dithering
func DitherRGB2PCT(
	red, green, blue, target RasterBand,
	ct ColorTable,
	progress ProgressFunc,
	data interface{},
) int {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	err := C.GDALDitherRGB2PCT(
		red.cval,
		green.cval,
		blue.cval,
		target.cval,
		ct.cval,
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	)
	return int(err)
}

// Compute checksum for image region
func (rb RasterBand) Checksum(xOff, yOff, xSize, ySize int) int {
	sum := C.GDALChecksumImage(rb.cval, C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize))
	return int(sum)
}

// Compute the proximity of all pixels in the image to a set of pixels in the source image
func (src RasterBand) ComputeProximity(
	dest RasterBand,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return CPLErr(
		C.GDALComputeProximity(
			src.cval,
			dest.cval,
			(**C.char)(unsafe.Pointer(&opts[0])),
			C.goGDALProgressFuncProxyB(),
			unsafe.Pointer(arg),
		),
	).Err()
}

// Fill selected raster regions by interpolation from the edges
func (src RasterBand) FillNoData(
	mask RasterBand,
	distance float64,
	iterations int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return CPLErr(
		C.GDALFillNodata(
			src.cval,
			mask.cval,
			C.double(distance),
			0,
			C.int(iterations),
			(**C.char)(unsafe.Pointer(&opts[0])),
			C.goGDALProgressFuncProxyB(),
			unsafe.Pointer(arg),
		),
	).Err()
}

// Create polygon coverage from raster data using an integer buffer
func (src RasterBand) Polygonize(
	mask RasterBand,
	layer Layer,
	fieldIndex int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return CPLErr(
		C.GDALPolygonize(
			src.cval,
			mask.cval,
			layer.cval,
			C.int(fieldIndex),
			(**C.char)(unsafe.Pointer(&opts[0])),
			C.goGDALProgressFuncProxyB(),
			unsafe.Pointer(arg),
		),
	).Err()
}

// Create polygon coverage from raster data using a floating point buffer
func (src RasterBand) FPolygonize(
	mask RasterBand,
	layer Layer,
	fieldIndex int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return CPLErr(
		C.GDALFPolygonize(
			src.cval,
			mask.cval,
			layer.cval,
			C.int(fieldIndex),
			(**C.char)(unsafe.Pointer(&opts[0])),
			C.goGDALProgressFuncProxyB(),
			unsafe.Pointer(arg),
		),
	).Err()
}

// Removes small raster polygons
func (src RasterBand) SieveFilter(
	mask, dest RasterBand,
	threshold, connectedness int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return CPLErr(
		C.GDALSieveFilter(
			src.cval,
			mask.cval,
			dest.cval,
			C.int(threshold),
			C.int(connectedness),
			(**C.char)(unsafe.Pointer(&opts[0])),
			C.goGDALProgressFuncProxyB(),
			unsafe.Pointer(arg),
		),
	).Err()
}

/* --------------------------------------------- */
/* Warp functions                                */
/* --------------------------------------------- */

//Unimplemented: CreateGenImgProjTransformer
//Unimplemented: CreateGenImgProjTransformer2
//Unimplemented: CreateGenImgProjTransformer3
//Unimplemented: SetGenImgProjTransformerDstGeoTransform
//Unimplemented: DestroyGenImgProjTransformer
//Unimplemented: GenImgProjTransform

//Unimplemented: CreateReprojectionTransformer
//Unimplemented: DestroyReprojection
//Unimplemented: ReprojectionTransform
//Unimplemented: CreateGCPTransformer
//Unimplemented: CreateGCPRefineTransformer
//Unimplemented: DestroyGCPTransformer
//Unimplemented: GCPTransform

//Unimplemented: CreateTPSTransformer
//Unimplemented: DestroyTPSTransformer
//Unimplemented: TPSTransform

type RPCInfoV2 struct {
	cval C.GDALRPCInfoV2
}

func ExtractRPCInfoV2(rpcMetadata []string) (RPCInfoV2, error) {
	// Convert rpcMetadata to C array
	opts := make([]*C.char, len(rpcMetadata)+1)
	for i, s := range rpcMetadata {
		opts[i] = C.CString(s)
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[len(rpcMetadata)] = (*C.char)(unsafe.Pointer(nil))

	var rpc RPCInfoV2
	passed := (*C.GDALRPCInfoV2)(unsafe.Pointer(&rpc.cval))
	if C.GDALExtractRPCInfoV2((**C.char)(unsafe.Pointer(&opts[0])), passed) == 0 {
		bigStr := strings.Join(rpcMetadata, "")
		if !strings.Contains(bigStr, RPC_LINE_NUM_COEFF) {
			return rpc, errors.New("RPC metadata does not contain LINE_NUM_COEFF")
		}
		if !strings.Contains(bigStr, RPC_LINE_DEN_COEFF) {
			return rpc, errors.New("RPC metadata does not contain LINE_DEN_COEFF")
		}
		if !strings.Contains(bigStr, RPC_SAMP_NUM_COEFF) {
			return rpc, errors.New("RPC metadata does not contain SAMP_NUM_COEFF")
		}
		if !strings.Contains(bigStr, RPC_SAMP_DEN_COEFF) {
			return rpc, errors.New("RPC metadata does not contain SAMP_DEN_COEFF")
		}

		return rpc, errors.New("GDALExtractRPCInfoV2 failed")
	}
	return rpc, nil
}

type RPCTransformer struct {
	// void*
	cval unsafe.Pointer
}

func CreateRPCTransformer(rpc RPCInfoV2, reversed bool, threshold float64, options []string) RPCTransformer {
	opts := make([]*C.char, len(options)+1)
	for i, s := range options {
		opts[i] = C.CString(s)
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[len(options)] = (*C.char)(unsafe.Pointer(nil))

	bReversed := 0
	if reversed {
		bReversed = 1
	}

	return RPCTransformer{C.GDALCreateRPCTransformerV2(
		&rpc.cval,
		C.int(bReversed),
		C.double(threshold),
		(**C.char)(unsafe.Pointer(&opts[0])),
	)}
}

func (t RPCTransformer) Destroy() {
	if t.cval != nil {
		C.GDALDestroyRPCTransformer(t.cval)
	}
}

// Transform a slice of x/y/z points using the RPC transformer. If reversed is true, the inverse transformation is applied.
// If the transformation fails for any point, an error is returned.
func (t RPCTransformer) Transform(x, y, z []float64, reversed bool) (xo, yo, zo []float64, err error) {
	if t.cval == nil {
		return nil, nil, nil, fmt.Errorf("RPCTransformer is not initialized")
	}
	nPoints := len(x)
	if nPoints != len(y) || nPoints != len(z) {
		return nil, nil, nil, fmt.Errorf("x, y, z slices must have the same length")
	}

	xo = make([]float64, nPoints)
	yo = make([]float64, nPoints)
	zo = make([]float64, nPoints)
	copy(xo, x)
	copy(yo, y)
	copy(zo, z)

	res := make([]int32, nPoints)
	C.GDALRPCTransform(
		t.cval,
		C.int(1),
		C.int(nPoints),
		(*C.double)(&xo[0]),
		(*C.double)(&yo[0]),
		(*C.double)(&zo[0]),
		(*C.int)(unsafe.Pointer(&res[0])),
	)

	for i, r := range res {
		if r == 0 {
			err = errors.Join(fmt.Errorf("rpc transform failed for (%f, %f, %f)", x[i], y[i], z[i]), err)
		}
	}

	return xo, yo, zo, nil
}

//Unimplemented: DestroyRPCTransformer
//Unimplemented: RPCTransform

//Unimplemented: CreateGeoLocTransformer
//Unimplemented: DestroyGeoLocTransformer
//Unimplemented: GeoLocTransform

//Unimplemented: CreateApproxTransformer
//Unimplemented: DestroyApproxTransformer
//Unimplemented: ApproxTransform

//Unimplemented: SimpleImageWarp
//Unimplemented: SuggestedWarpOutput
//Unimplemented: SuggestedWarpOutput2
//Unimplemented: SerializeTransformer
//Unimplemented: DeserializeTransformer

//Unimplemented: TransformGeolocations

/* --------------------------------------------- */
/* Contour line functions                        */
/* --------------------------------------------- */

//Unimplemented: CreateContourGenerator
//Unimplemented: FeedLine
//Unimplemented: Destroy
//Unimplemented: ContourWriter
//Unimplemented: ContourGenerate

/* --------------------------------------------- */
/* Rasterizer functions                          */
/* --------------------------------------------- */

// Burn geometries into raster
//Unimplmemented: RasterizeGeometries

// Burn geometries from the specified list of layers into the raster
//Unimplemented: RasterizeLayers

// Burn geometries from the specified list of layers into the raster
//Unimplemented: RasterizeLayersBuf

/* --------------------------------------------- */
/* Gridding functions                            */
/* --------------------------------------------- */

//Unimplemented: CreateGrid
//Unimplemented: ComputeMatchingPoints

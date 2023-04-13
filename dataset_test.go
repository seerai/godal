package gdal_test

import (
	"testing"

	gdal "github.com/seerai/godal"
	"github.com/stretchr/testify/assert"
)

var webMercatorWKT string = `
PROJCS["WGS 84 / Pseudo-Mercator",
    GEOGCS["WGS 84",
        DATUM["WGS_1984",
            SPHEROID["WGS 84",6378137,298.257223563,
                AUTHORITY["EPSG","7030"]],
            AUTHORITY["EPSG","6326"]],
        PRIMEM["Greenwich",0,
            AUTHORITY["EPSG","8901"]],
        UNIT["degree",0.0174532925199433,
            AUTHORITY["EPSG","9122"]],
        AUTHORITY["EPSG","4326"]],
    PROJECTION["Mercator_1SP"],
    PARAMETER["central_meridian",0],
    PARAMETER["scale_factor",1],
    PARAMETER["false_easting",0],
    PARAMETER["false_northing",0],
    UNIT["metre",1,
        AUTHORITY["EPSG","9001"]],
    AXIS["Easting",EAST],
    AXIS["Northing",NORTH],
    EXTENSION["PROJ4","+proj=merc +a=6378137 +b=6378137 +lat_ts=0 +lon_0=0 +x_0=0 +y_0=0 +k=1 +units=m +nadgrids=@null +wktext +no_defs"],
    AUTHORITY["EPSG","3857"]]`

func testDataset(t *testing.T) gdal.Dataset {
	driver, err := gdal.GetDriverByName("MEM")
	assert.NoError(t, err)

	ds := driver.Create("", 100, 100, 3, gdal.Float32, nil)
	ds.SetGeoTransform([6]float64{0.0, 1.0, 0.0, 1.0, 0.0, -1.0})
	ds.SetProjection(webMercatorWKT)

	return ds
}

func TestBasidReadWrite(t *testing.T) {

	ds := testDataset(t)

	writeBytes := make([]byte, ds.RasterCount()*ds.RasterXSize()*ds.RasterYSize()*4)
	for i := range writeBytes {
		writeBytes[i] = 1
	}

	err := ds.BasicWrite(0, 0, ds.RasterXSize(), ds.RasterYSize(), []int{1, 2, 3}, writeBytes)
	assert.NoError(t, err)

	readBytes := make([]byte, len(writeBytes))
	for i := range readBytes {
		readBytes[i] = 1
	}
	err = ds.BasicRead(0, 0, ds.RasterXSize(), ds.RasterYSize(), []int{1, 2, 3}, readBytes)
	assert.NoError(t, err)

	for i := range readBytes {
		assert.Equal(t, writeBytes[i], readBytes[i])
	}

}

func TestDatasetMetadata(t *testing.T) {
	d := gdal.Dataset{}
	md := d.Metadata("does not exist")
	assert.Equal(t, []string(nil), md)
}

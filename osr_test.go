package gdal_test

import (
	"fmt"
	"math"
	"testing"

	gdal "github.com/seerai/godal"
	"github.com/stretchr/testify/assert"
)

func TestToProjJSON(t *testing.T) {
	sr := gdal.CreateSpatialReference(nil)
	err := sr.FromEPSG(4326)
	assert.NoError(t, err)

	pj, err := sr.ToProjJSON()
	assert.NoError(t, err)
	fmt.Println(pj)
}

func TestLinearUnits(t *testing.T) {
	sr := gdal.CreateSpatialReference(nil)
	err := sr.FromURN("urn:ogc:def:crs:EPSG:6.18.3:3857")
	assert.NoError(t, err)

	units, value := sr.LinearUnits()
	assert.Equal(t, "metre", units)
	assert.Equal(t, 1.0, value)

	units, value = sr.TargetLinearUnits("PROJCS")
	assert.Equal(t, "metre", units)
	assert.Equal(t, 1.0, value)
	sr.Destroy()
}

func TestAngularUnits(t *testing.T) {
	sr := gdal.CreateSpatialReference(nil)
	err := sr.FromURN("urn:ogc:def:crs:EPSG::4326")
	assert.NoError(t, err)

	units, value := sr.AngularUnits()
	assert.Equal(t, "degree", units)
	assert.Equal(t, math.Pi/180, value)
	sr.Destroy()
}

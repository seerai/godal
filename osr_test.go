package gdal_test

import (
	"fmt"
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

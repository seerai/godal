package gdal_test

import (
	"testing"

	gdal "github.com/seerai/godal"
	"github.com/stretchr/testify/assert"
)

func TestMakeValid(t *testing.T) {
	// wkt is a self-intersecting polygon
	wkt := "POLYGON ((0 0, 1 1, 1 0, 0 1, 0 0))"
	g, err := gdal.CreateFromWKT(wkt, gdal.SpatialReference{})
	assert.NoError(t, err)
	assert.False(t, g.IsValid())

	valid, err := g.MakeValid()
	assert.NoError(t, err)
	assert.True(t, valid.IsValid())
	wktOut, err := valid.ToWKT()
	assert.NoError(t, err)
	assert.Equal(t, "MULTIPOLYGON (((0 1,0.5 0.5,0 0,0 1)),((1 0,0.5 0.5,1 1,1 0)))", wktOut)
}

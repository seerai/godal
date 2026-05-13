package gdal

import (
	"errors"
	"testing"
)

func TestOGRErrErr(t *testing.T) {
	cases := []struct {
		code int
		want error
	}{
		{0, nil},
		{1, ErrOGRNotEnoughData},
		{2, ErrOGRNotEnoughMemory},
		{3, ErrOGRUnsupportedGeometry},
		{4, ErrOGRUnsupportedOperation},
		{5, ErrOGRCorruptData},
		{6, ErrOGRFailure},
		{7, ErrOGRUnsupportedSRS},
		{8, ErrOGRInvalidHandle},
		{9, ErrOGRNonExistingFeature},
	}
	for _, c := range cases {
		got := OGRErr(c.code).Err()
		if c.want == nil {
			if got != nil {
				t.Errorf("OGRErr(%d).Err() = %v, want nil", c.code, got)
			}
			continue
		}
		if !errors.Is(got, c.want) {
			t.Errorf("OGRErr(%d).Err() = %v, want %v", c.code, got, c.want)
		}
	}

	got := OGRErr(99).Err()
	if !errors.Is(got, ErrIllegal) {
		t.Errorf("OGRErr(99).Err() = %v, want wrapped ErrIllegal", got)
	}
}

func TestCPLErrErr(t *testing.T) {
	cases := []struct {
		code int
		want error
	}{
		{0, nil},
		{1, ErrDebug},
		{2, ErrWarning},
		{3, ErrFailure},
		{4, ErrFatal},
	}
	for _, c := range cases {
		got := CPLErr(c.code).Err()
		if c.want == nil {
			if got != nil {
				t.Errorf("CPLErr(%d).Err() = %v, want nil", c.code, got)
			}
			continue
		}
		if !errors.Is(got, c.want) {
			t.Errorf("CPLErr(%d).Err() = %v, want %v", c.code, got, c.want)
		}
	}

	got := CPLErr(42).Err()
	if !errors.Is(got, ErrIllegal) {
		t.Errorf("CPLErr(42).Err() = %v, want wrapped ErrIllegal", got)
	}
}

// TestDataTypes tries various functions for gdal types
func TestDataTypes(t *testing.T) {
	var b DataType
	b = Byte

	if ret := b.Size(); ret != 8 {
		t.Fail()
	}

	if b.IsComplex() != 0 {
		t.Fail()
	}
}

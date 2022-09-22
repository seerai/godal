# GoDAL - Go Wrappers around the Geospatial Abstraction Library
[![Build](https://github.com/seerai/godal/actions/workflows/build_and_release.yaml/badge.svg)](https://github.com/seerai/godal/actions/workflows/build_and_release.yaml)
[![codecov](https://codecov.io/gh/seerai/godal/branch/main/graph/badge.svg?token=WQ4V3VYWHY)](https://codecov.io/gh/seerai/godal)

## About

The gdal.go package provides a go wrapper for GDAL, the Geospatial Data Abstraction Library. More information about GDAL can be found at http://www.gdal.org

This repo was originally forked from `github.com/lukeroth/gdal` and also `github.com/Rob-Fletcher` and will be developed further to closer mirror idiomatic Go.

It was also updated to support gdal 3.x, while the older version wouldn't work on later than 2.4. 
                                     
## Installation

1) go get github.com/seerai/godal
2) install libgdal-dev 3.5+
    - See [gdal's installation documentation for more details](https://gdal.org/download.html#). Note that many repos are not updated to 3.5 yet in which case there are [instructions for building from source](https://gdal.org/download.html#build-instructions).
3) go build 


## Compatibility

This software has been tested most recently on Ubuntu 20.04 with gdal 3.5

## Examples

See the examples directory for some examples of low level usage.  The macro package has higher level utilities

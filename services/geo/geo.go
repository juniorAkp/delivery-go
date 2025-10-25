package geo

import (
	"fmt"

	"github.com/juniorAkp/delivery-go/utils"
	"github.com/uber/h3-go/v4"
)

func CoordToIndex(coord utils.Coord) (h3.Cell, error) {
	latLng := h3.NewLatLng(coord.Longitude, coord.Latitude)
	resolution := 8

	return h3.LatLngToCell(latLng, resolution)
}

func GetNeighbors(cell h3.Cell, ringSize int) (neighbors []h3.Cell, err error) {
	if ringSize < 1 {
		return nil, fmt.Errorf("ring size must be at least 1")
	}

	neighbors, err = h3.GridDisk(cell, ringSize)
	if err != nil {
		return nil, err
	}
	return neighbors, nil
}

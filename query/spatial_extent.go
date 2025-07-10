package query

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
)

// SpatialExtent is an interface for types that can be converted to a spatial extent, to be used in queries.
type SpatialExtent interface {
	ToProtoSpatialFilter() (*datasetsv1.SpatialFilter, error)
}

type SpatialFilter struct {
	Geometry         orb.Geometry
	Mode             datasetsv1.SpatialFilterMode
	CoordinateSystem datasetsv1.SpatialCoordinateSystem
}

func (s *SpatialFilter) ToProtoSpatialFilter() (*datasetsv1.SpatialFilter, error) {
	switch s.Geometry.(type) {
	case orb.Polygon, orb.MultiPolygon:
		// ok
	default:
		return nil, fmt.Errorf("invalid geometry type, only Polygon and MultiPolygon are supported when querying: got %T", s.Geometry)
	}

	wkb, err := ewkb.Marshal(s.Geometry, ewkb.DefaultSRID)
	if err != nil {
		return nil, fmt.Errorf("invalid geometry: failed to marshal geometry as wkb: %w", err)
	}
	return &datasetsv1.SpatialFilter{
		Geometry:         &datasetsv1.Geometry{Wkb: wkb},
		Mode:             s.Mode,
		CoordinateSystem: s.CoordinateSystem,
	}, nil
}

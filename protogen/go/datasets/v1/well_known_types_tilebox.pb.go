// Code is NOT generated.

package datasetsv1

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// NewUUID constructs a new UUID from the provided uuid.UUID.
func NewUUID(id uuid.UUID) *UUID {
	return &UUID{Uuid: id[:]}
}

// AsUUID converts x to a uuid.UUID.
func (x *UUID) AsUUID() uuid.UUID {
	id, err := uuid.FromBytes(x.GetUuid())
	if err != nil {
		return uuid.Nil
	}
	return id
}

// IsValid reports whether the timestamp is valid.
// It is equivalent to CheckValid == nil.
func (x *UUID) IsValid() bool {
	return x.CheckValid() == nil
}

// CheckValid returns an error if the uuid is invalid.
// An error is reported for a nil UUID.
func (x *UUID) CheckValid() error {
	if x == nil {
		return protoimpl.X.NewError("invalid UUID: nil")
	}

	_, err := uuid.FromBytes(x.GetUuid())
	if err != nil {
		return protoimpl.X.NewError("invalid UUID: %v", err)
	}
	return nil
}

const sridWorldGeodetic1984 = 4326 // https://epsg.io/4326

// NewGeometry constructs a new Geometry from the provided orb.Geometry.
func NewGeometry(geometry orb.Geometry) *Geometry {
	wkb, err := ewkb.Marshal(geometry, sridWorldGeodetic1984)
	if err != nil {
		return nil
	}

	return &Geometry{Wkb: wkb}
}

// AsGeometry converts x to a orb.Geometry.
func (x *Geometry) AsGeometry() orb.Geometry {
	geometry, err := parseWKB(x.GetWkb())
	if err != nil {
		return nil
	}
	return geometry
}

// IsValid reports whether the geometry is valid.
// It is equivalent to CheckValid == nil.
func (x *Geometry) IsValid() bool {
	return x.CheckValid() == nil
}

// CheckValid returns an error if the geometry is invalid.
// An error is reported for a nil Geometry.
func (x *Geometry) CheckValid() error {
	if x == nil {
		return protoimpl.X.NewError("invalid geometry: nil")
	}

	_, err := parseWKB(x.GetWkb())
	if err != nil {
		return protoimpl.X.NewError("invalid geometry: %v", err)
	}
	return nil
}

// parseWKB parses a WKB geometry into an orb.Geometry. If the input is ewkb (extended wkb), additionally the SRID
// (spatial reference system identifier) is also parsed and validated to be either unset (0) or 4326 (WGS84).
func parseWKB(wkb []byte) (orb.Geometry, error) {
	if len(wkb) == 0 {
		return nil, errors.New("missing geometry")
	}

	geometry, srid, err := ewkb.Unmarshal(wkb)
	if err != nil {
		return nil, err
	}
	if srid != sridWorldGeodetic1984 && srid != 0 { // SRID 0 means not set, in which case we assume it is WGS84
		return nil, fmt.Errorf("expected geometry in EPSG 4326 (WGS84) SRID, got %d", srid)
	}
	return geometry, nil
}

package query

import (
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
)

func TestSpatialFilter_ToProtoSpatialFilter(t *testing.T) {
	colorado := orb.Polygon{{{-109.05, 41.00}, {-102.05, 41.00}, {-102.05, 37.0}, {-109.045, 37.0}, {-109.05, 41.00}}}
	coloradoWkb, err := ewkb.Marshal(colorado, ewkb.DefaultSRID)
	require.NoError(t, err)

	tests := []struct {
		name          string
		spatialFilter *SpatialFilter
		want          *datasetsv1.SpatialFilter
		wantErr       string
	}{
		{
			name: "ToProtoSpatialFilter",
			spatialFilter: &SpatialFilter{
				Geometry:         colorado,
				Mode:             datasetsv1.SpatialFilterMode_SPATIAL_FILTER_MODE_INTERSECTS,
				CoordinateSystem: datasetsv1.SpatialCoordinateSystem_SPATIAL_COORDINATE_SYSTEM_CARTESIAN,
			},
			want: &datasetsv1.SpatialFilter{
				Geometry:         &datasetsv1.Geometry{Wkb: coloradoWkb},
				Mode:             datasetsv1.SpatialFilterMode_SPATIAL_FILTER_MODE_INTERSECTS,
				CoordinateSystem: datasetsv1.SpatialCoordinateSystem_SPATIAL_COORDINATE_SYSTEM_CARTESIAN,
			},
		},
		{
			name: "ToProtoSpatialFilter invalid geometry",
			spatialFilter: &SpatialFilter{
				Geometry: orb.Point{0, 0},
			},
			wantErr: "invalid geometry type, only Polygon and MultiPolygon are supported when querying: got orb.Point",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.spatialFilter.ToProtoSpatialFilter()
			if tt.wantErr != "" {
				require.Error(t, err, "ToProtoSpatialFilter() error = %v, wantErr %v", err, tt.wantErr)
				assert.Contains(t, err.Error(), tt.wantErr, "error didn't contain expected message: '%s', got error '%s' instead.", tt.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.GetGeometry(), got.GetGeometry())
			assert.Equal(t, tt.want.GetMode(), got.GetMode())
			assert.Equal(t, tt.want.GetCoordinateSystem(), got.GetCoordinateSystem())
		})
	}
}

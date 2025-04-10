package datasetsv1

import (
	"encoding/hex"
	"testing"

	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGeometry(t *testing.T) {
	validGeometry := orb.Polygon{
		{{5, 3}, {5, 4}, {6, 4}, {6, 3}, {5, 3}},
	}
	validWkb, err := hex.DecodeString("0103000020e610000001000000050000000000000000001440000000000000084000000000000014400000000000001040000000000000184000000000000010400000000000001840000000000000084000000000000014400000000000000840")
	require.NoError(t, err)
	geometryProto := &Geometry{Wkb: validWkb}

	t.Run("NewGeometry", func(t *testing.T) {
		got := NewGeometry(validGeometry)
		assert.Equal(t, geometryProto, got)
		assert.Equal(t, validGeometry, got.AsGeometry())
	})

	t.Run("AsGeometry", func(t *testing.T) {
		got := geometryProto.AsGeometry()
		assert.Equal(t, validGeometry, got)
		assert.Equal(t, geometryProto, NewGeometry(got))
	})
}

package datasets

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/query"
	"pgregory.net/rapid"
)

func TestCollection_String(t *testing.T) {
	genTimeInterval := rapid.OneOf(
		rapid.Just(query.NewEmptyTimeInterval()),
		rapid.Ptr(rapid.Custom(func(t *rapid.T) query.TimeInterval {
			return query.TimeInterval{
				Start:          time.Unix(rapid.Int64().Draw(t, "Start"), 0).UTC(),
				End:            time.Unix(rapid.Int64().Draw(t, "End"), 0).UTC(),
				StartExclusive: rapid.Bool().Draw(t, "StartExclusive"),
				EndInclusive:   rapid.Bool().Draw(t, "EndInclusive"),
			}
		}), true),
	)

	genCollection := rapid.Custom(func(t *rapid.T) *Collection {
		return &Collection{
			ID:           uuid.New(),
			Name:         rapid.String().Draw(t, "Name"),
			Availability: genTimeInterval.Draw(t, "Availability"),
			Count:        rapid.Uint64().Draw(t, "Count"),
		}
	})

	rapid.Check(t, func(t *rapid.T) {
		input := genCollection.Draw(t, "collection")
		got := input.String()

		assert.Contains(t, got, input.Name)
		assert.NotContains(t, got, input.ID.String())

		if input.Availability == nil {
			assert.Contains(t, got, "unknown")
		} else {
			assert.Contains(t, got, input.Availability.String())
		}

		if input.Count == 0 {
			assert.NotContains(t, got, "data points")
		} else {
			assert.Contains(t, got, fmt.Sprintf("%d data points", input.Count))
		}
	})
}

func Test_collectionClient_Get(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "collection")

	dataset, err := client.Datasets.Get(ctx, "open_data.asf.ers_sar")
	require.NoError(t, err)

	collection, err := client.Collections.Get(ctx, dataset.ID, "ERS-2")
	require.NoError(t, err)

	assert.Equal(t, "ERS-2", collection.Name)
	assert.Equal(t, "18b986bf-de32-4a76-ad20-39b3df056803", collection.ID.String())
	assert.Equal(t, "1995-10-01 03:13:03 +0000 UTC", collection.Availability.Start.String())
	assert.NotZero(t, collection.Count)
}

// testing only the Get part
func Test_collectionClient_GetOrCreate(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "collection_get_or_create")

	dataset, err := client.Datasets.Get(ctx, "open_data.asf.ers_sar")
	require.NoError(t, err)

	collection, err := client.Collections.GetOrCreate(ctx, dataset.ID, "ERS-2")
	require.NoError(t, err)

	assert.Equal(t, "ERS-2", collection.Name)
	assert.Equal(t, "18b986bf-de32-4a76-ad20-39b3df056803", collection.ID.String())
	assert.Equal(t, "1995-10-01 03:13:03 +0000 UTC", collection.Availability.Start.String())
	assert.NotZero(t, collection.Count)
}

func Test_collectionClient_List(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "collections")

	dataset, err := client.Datasets.Get(ctx, "open_data.asf.ers_sar")
	require.NoError(t, err)

	collections, err := client.Collections.List(ctx, dataset.ID)
	require.NoError(t, err)

	names := lo.Map(collections, func(c *Collection, _ int) string {
		return c.Name
	})
	assert.Equal(t, []string{"ERS-1", "ERS-2"}, names)
}

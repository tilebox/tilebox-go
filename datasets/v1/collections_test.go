package datasets

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_collectionClient_Get(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "collection")

	dataset, err := client.Datasets.Get(ctx, "open_data.asf.ers_sar")
	require.NoError(t, err)

	collection, err := client.Collections.Get(ctx, dataset.ID, "ERS-2")
	require.NoError(t, err)

	assert.Equal(t, "ERS-2", collection.Name)
	assert.Equal(t, "c408f2b8-0488-4528-9fb7-a18361df639f", collection.ID.String())
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

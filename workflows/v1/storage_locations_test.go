package workflows

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/client"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/query"
	"github.com/tilebox/tilebox-go/workflows/v1/storagelocation"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestStorageLocationReferences(t *testing.T) {
	tests := []struct {
		name      string
		reference *StorageLocationReference
		wire      *workflowsv1.StorageLocationReference
	}{
		{
			name:      "aws-s3",
			reference: &StorageLocationReference{Type: StorageTypeAWSS3, AWSS3Bucket: &AWSS3BucketReference{Bucket: "imagery", Region: "eu-west-1"}},
			wire: workflowsv1.StorageLocationReference_builder{
				Type:        workflowsv1.StorageType_STORAGE_TYPE_AWS_S3,
				AwsS3Bucket: workflowsv1.AWSS3BucketReference_builder{Bucket: "imagery", Region: "eu-west-1"}.Build(),
			}.Build(),
		},
		{
			name:      "gcs",
			reference: &StorageLocationReference{Type: StorageTypeGCS, GCSBucket: &GCSBucketReference{Bucket: "scenes", ProjectID: "my-project", Location: "EUR4"}},
			wire: workflowsv1.StorageLocationReference_builder{
				Type:      workflowsv1.StorageType_STORAGE_TYPE_GCS,
				GcsBucket: workflowsv1.GCSBucketReference_builder{Bucket: "scenes", ProjectId: "my-project", Location: "EUR4"}.Build(),
			}.Build(),
		},
		{
			name: "azure-blob",
			reference: &StorageLocationReference{Type: StorageTypeAzureBlob, AzureBlob: &AzureBlobReference{
				StorageAccountResourceID: "/subscriptions/sub/resourceGroups/group/providers/Microsoft.Storage/storageAccounts/account", Container: "tiles", Region: "westeurope",
			}},
			wire: workflowsv1.StorageLocationReference_builder{
				Type: workflowsv1.StorageType_STORAGE_TYPE_AZURE_BLOB,
				AzureBlob: workflowsv1.AzureBlobReference_builder{
					StorageAccountResourceId: "/subscriptions/sub/resourceGroups/group/providers/Microsoft.Storage/storageAccounts/account", Container: "tiles", Region: "westeurope",
				}.Build(),
			}.Build(),
		},
		{
			name:      "local",
			reference: &StorageLocationReference{Type: StorageTypeLocal, Local: &LocalReference{Path: "/data/incoming"}},
			wire: workflowsv1.StorageLocationReference_builder{
				Type:  workflowsv1.StorageType_STORAGE_TYPE_LOCAL,
				Local: workflowsv1.LocalReference_builder{Path: "/data/incoming"}.Build(),
			}.Build(),
		},
		{name: "nil"},
		{name: "unspecified", reference: &StorageLocationReference{Type: StorageTypeUnspecified}, wire: &workflowsv1.StorageLocationReference{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, proto.Equal(tt.wire, tt.reference.toProto()))
			assert.Equal(t, tt.reference, protoToStorageLocationReference(tt.wire))
			if tt.reference != nil {
				assert.Equal(t, tt.name, tt.reference.Type.String())
				encoded, err := json.Marshal(tt.reference.Type)
				require.NoError(t, err)
				assert.JSONEq(t, strconv.Quote(tt.name), string(encoded))
			}
		})
	}
	assert.Equal(t, StorageTypeUnspecified, protoToStorageType(workflowsv1.StorageType(99)))
	assert.Nil(t, protoToStorageLocation(nil))
	assert.Nil(t, protoToStorageLocation(&workflowsv1.StorageLocation{}).Reference)
}

func TestStorageLocationClientRPCs(t *testing.T) {
	locationID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1546")
	subscriptionID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1547")
	cursorID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1548")
	reference := workflowsv1.StorageLocationReference_builder{
		Type:        workflowsv1.StorageType_STORAGE_TYPE_AWS_S3,
		AwsS3Bucket: workflowsv1.AWSS3BucketReference_builder{Bucket: "input", Region: "eu-central-1"}.Build(),
	}.Build()
	wireLocation := workflowsv1.StorageLocation_builder{Id: tileboxv1.NewUUID(locationID), Name: "Input", Reference: reference}.Build()
	location := &StorageLocation{ID: locationID, Name: "Input", Reference: &StorageLocationReference{
		Type: StorageTypeAWSS3, AWSS3Bucket: &AWSS3BucketReference{Bucket: "input", Region: "eu-central-1"},
	}}
	wireSubscription := workflowsv1.StorageSubscription_builder{
		Id: tileboxv1.NewUUID(subscriptionID), StorageLocationId: tileboxv1.NewUUID(locationID),
		Type: workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_AWS_SNS, Endpoint: "https://example.com/notify",
		AwsSns: workflowsv1.AWSSNSStorageSubscription_builder{TopicArn: "arn:aws:sns:eu-central-1:123456789012:objects", MessageFormat: "earthsearch_stac"}.Build(),
	}.Build()
	subscription := &StorageSubscription{
		ID: subscriptionID, StorageLocationID: locationID, Type: StorageSubscriptionTypeAWSSNS, Endpoint: "https://example.com/notify",
		AWSSNS: &AWSSNSStorageSubscription{TopicARN: "arn:aws:sns:eu-central-1:123456789012:objects", MessageFormat: "earthsearch_stac"},
	}
	tests := []struct {
		method   string
		request  proto.Message
		response proto.Message
		call     func(StorageLocationClient) (any, error)
		want     any
	}{
		{
			method:  "CreateStorageLocation",
			request: workflowsv1.CreateStorageLocationRequest_builder{Name: "Input", Reference: reference}.Build(), response: wireLocation,
			call: func(c StorageLocationClient) (any, error) { return c.Create(t.Context(), "Input", location.Reference) }, want: location,
		},
		{
			method: "GetStorageLocation", request: tileboxv1.NewUUID(locationID), response: wireLocation,
			call: func(c StorageLocationClient) (any, error) { return c.Get(t.Context(), locationID) }, want: location,
		},
		{
			method: "ListStorageLocations", request: &emptypb.Empty{},
			response: workflowsv1.StorageLocations_builder{Locations: []*workflowsv1.StorageLocation{wireLocation}}.Build(),
			call:     func(c StorageLocationClient) (any, error) { return c.List(t.Context()) }, want: []*StorageLocation{location},
		},
		{
			method:  "UpdateStorageLocation",
			request: workflowsv1.UpdateStorageLocationRequest_builder{StorageLocationId: tileboxv1.NewUUID(locationID), Name: "Input"}.Build(), response: wireLocation,
			call: func(c StorageLocationClient) (any, error) { return c.Update(t.Context(), locationID, "Input") }, want: location,
		},
		{
			method: "DeleteStorageLocation", request: tileboxv1.NewUUID(locationID), response: &emptypb.Empty{},
			call: func(c StorageLocationClient) (any, error) { return nil, c.Delete(t.Context(), locationID) },
		},
		{
			method: "CreateStorageSubscription",
			request: workflowsv1.CreateStorageSubscriptionRequest_builder{
				StorageLocationId: tileboxv1.NewUUID(locationID), Type: workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_AWS_SNS,
				AwsSns: workflowsv1.AWSSNSStorageSubscription_builder{TopicArn: "arn:aws:sns:eu-central-1:123456789012:objects", MessageFormat: "earthsearch_stac"}.Build(),
			}.Build(), response: wireSubscription,
			call: func(c StorageLocationClient) (any, error) {
				return c.CreateSubscription(t.Context(), locationID, StorageSubscriptionConfig{Type: StorageSubscriptionTypeAWSSNS, AWSSNS: subscription.AWSSNS})
			}, want: subscription,
		},
		{
			method: "GetStorageSubscription", request: tileboxv1.NewUUID(subscriptionID), response: wireSubscription,
			call: func(c StorageLocationClient) (any, error) { return c.GetSubscription(t.Context(), subscriptionID) }, want: subscription,
		},
		{
			method:   "ListStorageSubscriptions",
			request:  workflowsv1.ListStorageSubscriptionsRequest_builder{StorageLocationId: tileboxv1.NewUUID(locationID)}.Build(),
			response: workflowsv1.StorageSubscriptions_builder{Subscriptions: []*workflowsv1.StorageSubscription{wireSubscription}}.Build(),
			call:     func(c StorageLocationClient) (any, error) { return c.ListSubscriptions(t.Context(), locationID) }, want: []*StorageSubscription{subscription},
		},
		{
			method: "DeleteStorageSubscription", request: tileboxv1.NewUUID(subscriptionID), response: &emptypb.Empty{},
			call: func(c StorageLocationClient) (any, error) {
				return nil, c.DeleteSubscription(t.Context(), subscriptionID)
			},
		},
		{
			method: "ListStorageSubscriptionEvents",
			request: workflowsv1.ListStorageSubscriptionEventsRequest_builder{
				StorageLocationId: tileboxv1.NewUUID(locationID), Page: tileboxv1.Pagination_builder{Limit: proto.Int64(7), StartingAfter: tileboxv1.NewUUID(cursorID)}.Build(),
			}.Build(),
			response: workflowsv1.StorageSubscriptionEvents_builder{
				Events: []*workflowsv1.StorageSubscriptionEvent{workflowsv1.StorageSubscriptionEvent_builder{
					Id: tileboxv1.NewUUID(cursorID), StorageSubscriptionId: tileboxv1.NewUUID(subscriptionID), ObjectKey: "a/b.tif", Type: workflowsv1.StorageEventType_STORAGE_EVENT_TYPE_CREATED,
				}.Build()},
				NextPage: tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(subscriptionID)}.Build(),
			}.Build(),
			call: func(c StorageLocationClient) (any, error) {
				return c.ListSubscriptionEventsPage(t.Context(), locationID, storagelocation.WithLimit(7), storagelocation.WithCursor(query.NewCursor(cursorID)))
			},
			want: &StorageSubscriptionEventPage{
				Events:     []*StorageSubscriptionEvent{{ID: cursorID, StorageSubscriptionID: subscriptionID, ObjectKey: "a/b.tif", Type: StorageEventTypeCreated, TriggeredJobs: []*TriggeredJob{}}},
				NextCursor: query.NewCursor(subscriptionID),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			var fail atomic.Bool
			response, err := proto.Marshal(tt.response)
			require.NoError(t, err)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/workflows.v1.StorageLocationService/"+tt.method, r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
				assert.NotEmpty(t, r.Header.Get("Tilebox-Client"))
				body, readErr := io.ReadAll(r.Body)
				if !assert.NoError(t, readErr) {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				request := tt.request.ProtoReflect().New().Interface()
				assert.NoError(t, proto.Unmarshal(body, request))
				assert.True(t, proto.Equal(tt.request, request), "request: %v", request)
				if fail.Load() {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusNotFound)
					_, writeErr := io.WriteString(w, `{"code":"not_found","message":"missing"}`)
					assert.NoError(t, writeErr)
					return
				}
				w.Header().Set("Content-Type", "application/proto")
				_, writeErr := w.Write(response)
				assert.NoError(t, writeErr)
			}))
			t.Cleanup(server.Close)
			c := NewClient(WithURL(server.URL), WithHTTPClient(server.Client()), WithAPIKey("test-key"),
				WithClientMetadata(client.Metadata{Name: "test", Version: "1"}), WithDisableTracing(), WithDisableLogging())
			result, err := tt.call(c.StorageLocations)
			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
			fail.Store(true)
			_, err = tt.call(c.StorageLocations)
			require.Error(t, err)
			assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
			assert.Contains(t, err.Error(), "failed to ")
		})
	}
}

type storageEventsService struct {
	_storageLocationService

	list func(context.Context, uuid.UUID, *tileboxv1.Pagination) (*workflowsv1.StorageSubscriptionEvents, error)
}

func (s storageEventsService) ListStorageSubscriptionEvents(ctx context.Context, id uuid.UUID, page *tileboxv1.Pagination) (*workflowsv1.StorageSubscriptionEvents, error) {
	return s.list(ctx, id, page)
}

func TestStorageSubscriptionEventsPagination(t *testing.T) {
	locationID := uuid.New()
	startID := uuid.New()
	secondID := uuid.New()
	thirdID := uuid.New()
	for _, limit := range []int64{0, -1, 1, 3, 100, 101} {
		t.Run(strconv.FormatInt(limit, 10), func(t *testing.T) {
			calls := 0
			c := storageLocationClient{service: storageEventsService{list: func(_ context.Context, id uuid.UUID, page *tileboxv1.Pagination) (*workflowsv1.StorageSubscriptionEvents, error) {
				assert.Equal(t, locationID, id)
				calls++
				if calls == 1 {
					assert.Equal(t, startID, page.GetStartingAfter().AsUUID())
					assert.Equal(t, limit > 0, page.HasLimit())
					if limit > 0 {
						assert.Equal(t, min(limit, 100), page.GetLimit())
					}
					return workflowsv1.StorageSubscriptionEvents_builder{
						Events: []*workflowsv1.StorageSubscriptionEvent{
							workflowsv1.StorageSubscriptionEvent_builder{ObjectKey: "first"}.Build(),
							workflowsv1.StorageSubscriptionEvent_builder{ObjectKey: "second"}.Build(),
						},
						NextPage: tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(secondID)}.Build(),
					}.Build(), nil
				}
				if calls == 2 {
					assert.Equal(t, secondID, page.GetStartingAfter().AsUUID())
					if limit > 0 {
						assert.Equal(t, limit-2, page.GetLimit())
					}
					// An empty page with a continuation must not terminate iteration.
					return workflowsv1.StorageSubscriptionEvents_builder{NextPage: tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(thirdID)}.Build()}.Build(), nil
				}
				require.Equal(t, 3, calls)
				assert.Equal(t, thirdID, page.GetStartingAfter().AsUUID())
				return workflowsv1.StorageSubscriptionEvents_builder{Events: []*workflowsv1.StorageSubscriptionEvent{
					workflowsv1.StorageSubscriptionEvent_builder{ObjectKey: "third"}.Build(),
					workflowsv1.StorageSubscriptionEvent_builder{ObjectKey: "fourth"}.Build(),
				}}.Build(), nil
			}}}
			iterator := c.ListSubscriptionEvents(t.Context(), locationID, storagelocation.WithCursor(query.NewCursor(startID)), storagelocation.WithLimit(limit))
			assert.Zero(t, calls, "iterator must be lazy")
			for range 2 {
				calls = 0
				var keys []string
				for event, err := range iterator {
					require.NoError(t, err)
					keys = append(keys, event.ObjectKey)
				}
				want := []string{"first", "second", "third", "fourth"}
				if limit > 0 && limit < 4 {
					want = want[:limit]
				}
				assert.Equal(t, want, keys)
				if limit == 1 {
					assert.Equal(t, 1, calls)
				} else {
					assert.Equal(t, 3, calls)
				}
			}
		})
	}
}

func TestStorageSubscriptionEventsStopAndError(t *testing.T) {
	calls := 0
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	c := storageLocationClient{service: storageEventsService{list: func(ctx context.Context, _ uuid.UUID, _ *tileboxv1.Pagination) (*workflowsv1.StorageSubscriptionEvents, error) {
		calls++
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return workflowsv1.StorageSubscriptionEvents_builder{
			Events:   []*workflowsv1.StorageSubscriptionEvent{workflowsv1.StorageSubscriptionEvent_builder{ObjectKey: "first"}.Build()},
			NextPage: tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(uuid.New())}.Build(),
		}.Build(), nil
	}}}
	for event, err := range c.ListSubscriptionEvents(ctx, uuid.New()) {
		require.NoError(t, err)
		assert.Equal(t, "first", event.ObjectKey)
		break
	}
	assert.Equal(t, 1, calls)
	calls = 0
	yields := 0
	for event, err := range c.ListSubscriptionEvents(ctx, uuid.New()) {
		yields++
		if yields == 1 {
			require.NoError(t, err)
			cancel()
		} else {
			require.ErrorIs(t, err, context.Canceled)
			assert.Nil(t, event)
		}
	}
	assert.Equal(t, 2, yields)
	assert.Equal(t, 2, calls)
}

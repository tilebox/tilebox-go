package workflows

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestStorageSubscriptionProviders(t *testing.T) {
	locationID := uuid.New()
	subscriptionID := uuid.New()
	webhookSecret := uuid.NewString()
	createdAt := time.Date(2026, 10, 5, 12, 13, 14, 0, time.UTC)
	tests := []struct {
		name     string
		config   StorageSubscriptionConfig
		request  *workflowsv1.CreateStorageSubscriptionRequest
		response *workflowsv1.StorageSubscription
		want     *StorageSubscription
	}{
		{
			name: "sns",
			config: StorageSubscriptionConfig{Type: StorageSubscriptionTypeAWSSNS, AWSSNS: &AWSSNSStorageSubscription{
				TopicARN: "arn:aws:sns:us-west-2:123456789012:landsat", MessageFormat: "usgs_landsat_c2",
			}},
			request: workflowsv1.CreateStorageSubscriptionRequest_builder{
				StorageLocationId: tileboxv1.NewUUID(locationID), Type: workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_AWS_SNS,
				AwsSns: workflowsv1.AWSSNSStorageSubscription_builder{TopicArn: "arn:aws:sns:us-west-2:123456789012:landsat", MessageFormat: "usgs_landsat_c2"}.Build(),
			}.Build(),
			response: workflowsv1.StorageSubscription_builder{
				Type:   workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_AWS_SNS,
				AwsSns: workflowsv1.AWSSNSStorageSubscription_builder{TopicArn: "arn:aws:sns:us-west-2:123456789012:landsat", MessageFormat: "usgs_landsat_c2"}.Build(),
			}.Build(),
			want: &StorageSubscription{Type: StorageSubscriptionTypeAWSSNS, AWSSNS: &AWSSNSStorageSubscription{
				TopicARN: "arn:aws:sns:us-west-2:123456789012:landsat", MessageFormat: "usgs_landsat_c2",
			}},
		},
		{
			name: "pubsub",
			config: StorageSubscriptionConfig{Type: StorageSubscriptionTypeGooglePubSub, GooglePubSub: &GooglePubSubStorageSubscription{
				Subscription: "projects/project/subscriptions/objects", ServiceAccountEmail: "push@project.iam.gserviceaccount.com", Audience: "output-only",
			}},
			request: workflowsv1.CreateStorageSubscriptionRequest_builder{
				StorageLocationId: tileboxv1.NewUUID(locationID), Type: workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_GOOGLE_PUBSUB,
				GooglePubsub: workflowsv1.GooglePubSubStorageSubscription_builder{
					Subscription: "projects/project/subscriptions/objects", ServiceAccountEmail: "push@project.iam.gserviceaccount.com",
				}.Build(),
			}.Build(),
			response: workflowsv1.StorageSubscription_builder{
				Type: workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_GOOGLE_PUBSUB,
				GooglePubsub: workflowsv1.GooglePubSubStorageSubscription_builder{
					Subscription: "projects/project/subscriptions/objects", ServiceAccountEmail: "push@project.iam.gserviceaccount.com", Audience: "https://example.com/notify",
				}.Build(),
			}.Build(),
			want: &StorageSubscription{Type: StorageSubscriptionTypeGooglePubSub, GooglePubSub: &GooglePubSubStorageSubscription{
				Subscription: "projects/project/subscriptions/objects", ServiceAccountEmail: "push@project.iam.gserviceaccount.com", Audience: "https://example.com/notify",
			}},
		},
		{
			name:   "event grid",
			config: StorageSubscriptionConfig{Type: StorageSubscriptionTypeAzureEventGrid},
			request: workflowsv1.CreateStorageSubscriptionRequest_builder{
				StorageLocationId: tileboxv1.NewUUID(locationID), Type: workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_AZURE_EVENT_GRID,
				AzureEventGrid: &workflowsv1.AzureEventGridStorageSubscription{},
			}.Build(),
			response: workflowsv1.StorageSubscription_builder{
				Type:           workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_AZURE_EVENT_GRID,
				AzureEventGrid: workflowsv1.AzureEventGridStorageSubscription_builder{WebhookSecretHeader: "X-Tilebox-Webhook-Secret", WebhookSecret: webhookSecret}.Build(), //nolint:gosec // Public header name; the test secret is generated.
			}.Build(),
			want: &StorageSubscription{Type: StorageSubscriptionTypeAzureEventGrid, AzureEventGrid: &AzureEventGridStorageSubscription{ //nolint:gosec // Public header name; the test secret is generated.
				WebhookSecretHeader: "X-Tilebox-Webhook-Secret", WebhookSecret: webhookSecret,
			}},
		},
		{
			name:   "tilebox cli",
			config: StorageSubscriptionConfig{Type: StorageSubscriptionTypeTileboxCLI},
			request: workflowsv1.CreateStorageSubscriptionRequest_builder{
				StorageLocationId: tileboxv1.NewUUID(locationID), Type: workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_TILEBOX_CLI,
			}.Build(),
			response: workflowsv1.StorageSubscription_builder{Type: workflowsv1.StorageSubscriptionType_STORAGE_SUBSCRIPTION_TYPE_TILEBOX_CLI}.Build(),
			want:     &StorageSubscription{Type: StorageSubscriptionTypeTileboxCLI},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, proto.Equal(tt.request, tt.config.toProto(locationID)))
			tt.response.SetId(tileboxv1.NewUUID(subscriptionID))
			tt.response.SetStorageLocationId(tileboxv1.NewUUID(locationID))
			tt.response.SetCreatedAt(timestamppb.New(createdAt))
			tt.response.SetEndpoint("https://example.com/notify")
			tt.want.ID = subscriptionID
			tt.want.StorageLocationID = locationID
			tt.want.CreatedAt = createdAt
			tt.want.Endpoint = "https://example.com/notify"
			assert.Equal(t, tt.want, protoToStorageSubscription(tt.response))
		})
	}
	assert.Nil(t, protoToStorageSubscription(nil))
	assert.Equal(t, StorageSubscriptionTypeUnspecified, protoToStorageSubscriptionType(workflowsv1.StorageSubscriptionType(99)))
	assert.True(t, protoToStorageSubscription(&workflowsv1.StorageSubscription{}).CreatedAt.IsZero())
}

func TestStorageSubscriptionEventConversion(t *testing.T) {
	eventID, subscriptionID, jobID, automationID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	eventTime := time.Date(2026, 10, 5, 10, 11, 12, 123, time.UTC)
	receivedAt := eventTime.Add(5 * time.Second)
	wire := workflowsv1.StorageSubscriptionEvent_builder{
		Id: tileboxv1.NewUUID(eventID), StorageSubscriptionId: tileboxv1.NewUUID(subscriptionID), ObjectKey: "nested/object.tif",
		Type: workflowsv1.StorageEventType_STORAGE_EVENT_TYPE_CREATED, EventTime: timestamppb.New(eventTime), ReceivedAt: timestamppb.New(receivedAt),
		TriggeredJobs: []*workflowsv1.TriggeredJob{workflowsv1.TriggeredJob_builder{
			JobId: tileboxv1.NewUUID(jobID), AutomationId: tileboxv1.NewUUID(automationID), GlobPattern: "nested/*.tif",
		}.Build()},
	}.Build()
	assert.Equal(t, &StorageSubscriptionEvent{
		ID: eventID, StorageSubscriptionID: subscriptionID, ObjectKey: "nested/object.tif", Type: StorageEventTypeCreated,
		EventTime: &eventTime, ReceivedAt: receivedAt,
		TriggeredJobs: []*TriggeredJob{{JobID: jobID, AutomationID: automationID, GlobPattern: "nested/*.tif"}},
	}, protoToStorageSubscriptionEvent(wire))
	wire.ClearEventTime()
	require.Nil(t, protoToStorageSubscriptionEvent(wire).EventTime)
	assert.Nil(t, protoToStorageSubscriptionEvent(nil))
	assert.Equal(t, StorageEventTypeUnspecified, protoToStorageEventType(workflowsv1.StorageEventType(99)))
	assert.True(t, protoToStorageSubscriptionEvent(&workflowsv1.StorageSubscriptionEvent{}).ReceivedAt.IsZero())
}

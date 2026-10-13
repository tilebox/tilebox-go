package workflows

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/client"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/workflows/v1/workflowsv1connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAutomationStorageLocationReference(t *testing.T) {
	storageLocationID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1546")
	automation := protoToAutomation(workflowsv1.AutomationPrototype_builder{
		StorageEventTriggers: []*workflowsv1.StorageEventTrigger{
			workflowsv1.StorageEventTrigger_builder{
				GlobPattern: "data/*.tif",
				StorageLocation: workflowsv1.StorageLocation_builder{
					Id: tileboxv1.NewUUID(storageLocationID), Name: "Imagery",
					Reference: workflowsv1.StorageLocationReference_builder{
						Type:      workflowsv1.StorageType_STORAGE_TYPE_GCS,
						GcsBucket: workflowsv1.GCSBucketReference_builder{Bucket: "bucket", ProjectId: "project", Location: "EU"}.Build(),
					}.Build(),
				}.Build(),
			}.Build(),
		},
	}.Build())
	require.Len(t, automation.StorageEventTriggers, 1)
	assert.Equal(t, "data/*.tif", automation.StorageEventTriggers[0].GlobPattern)
	location := automation.StorageEventTriggers[0].StorageLocation
	assert.Equal(t, storageLocationID, location.ID)
	assert.Equal(t, "Imagery", location.Name)
	assert.Equal(t, &StorageLocationReference{
		Type: StorageTypeGCS, GCSBucket: &GCSBucketReference{Bucket: "bucket", ProjectID: "project", Location: "EU"},
	}, location.Reference)
}

func TestAutomationClientWrites(t *testing.T) {
	automationID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1546")
	triggerID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1547")
	locationID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1548")
	secondLocationID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1549")
	for _, method := range []string{"CreateAutomation", "UpdateAutomation"} {
		for _, kind := range []string{"cron", "storage"} {
			t.Run(method+"/"+kind, func(t *testing.T) {
				id, retainedTriggerID := uuid.Nil, uuid.Nil
				if method == "UpdateAutomation" {
					id, retainedTriggerID = automationID, triggerID
				}
				input := &Automation{
					ID: id, Name: "Process arrivals", Disabled: method == "CreateAutomation",
					Prototype: &AutomationTaskPrototype{
						ClusterSlug: "processing", Identifier: NewTaskIdentifier("Process", "v2.3"), Display: "Process scene",
						MaxRetries: 4, Input: []byte{0x0a, 0x03, 'a', 'b', 'c'},
					},
					CronTriggers: []*CronTrigger{}, StorageEventTriggers: []*StorageEventTrigger{},
				}
				request := workflowsv1.AutomationPrototype_builder{
					Id: tileboxv1.NewUUID(id), Name: "Process arrivals", Disabled: method == "CreateAutomation",
					Prototype: workflowsv1.SingleTaskSubmission_builder{
						ClusterSlug: "processing", Identifier: workflowsv1.TaskIdentifier_builder{Name: "Process", Version: "v2.3"}.Build(),
						Display: "Process scene", MaxRetries: 4, Input: []byte{0x0a, 0x03, 'a', 'b', 'c'},
					}.Build(),
				}.Build()
				if kind == "cron" {
					input.CronTriggers = []*CronTrigger{
						{ID: retainedTriggerID, Schedule: "0 6 * * 1-5"}, {Schedule: "CRON_TZ=Europe/Vienna 30 9 * * 6"},
					}
					request.SetCronTriggers([]*workflowsv1.CronTrigger{
						workflowsv1.CronTrigger_builder{Id: tileboxv1.NewUUID(retainedTriggerID), Schedule: "0 6 * * 1-5"}.Build(),
						workflowsv1.CronTrigger_builder{Id: tileboxv1.NewUUID(uuid.Nil), Schedule: "CRON_TZ=Europe/Vienna 30 9 * * 6"}.Build(),
					})
				} else {
					input.StorageEventTriggers = []*StorageEventTrigger{
						{ID: retainedTriggerID, StorageLocation: &StorageLocation{ID: locationID, Name: "Read-only location name"}, GlobPattern: "scenes/**/*.tif"},
						{StorageLocation: &StorageLocation{ID: secondLocationID}, GlobPattern: "metadata/*.json"},
					}
					request.SetStorageEventTriggers([]*workflowsv1.StorageEventTrigger{
						workflowsv1.StorageEventTrigger_builder{
							Id: tileboxv1.NewUUID(retainedTriggerID), StorageLocation: workflowsv1.StorageLocation_builder{Id: tileboxv1.NewUUID(locationID)}.Build(), GlobPattern: "scenes/**/*.tif",
						}.Build(),
						workflowsv1.StorageEventTrigger_builder{
							Id: tileboxv1.NewUUID(uuid.Nil), StorageLocation: workflowsv1.StorageLocation_builder{Id: tileboxv1.NewUUID(secondLocationID)}.Build(), GlobPattern: "metadata/*.json",
						}.Build(),
					})
				}
				response := proto.CloneOf(request)
				response.SetId(tileboxv1.NewUUID(automationID))
				response.SetName("Server response")
				var fail atomic.Bool
				c := newAutomationTestClient(t, method, automationRPCService{write: func(_ context.Context, got *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error) {
					assert.True(t, proto.Equal(request, got), "request: %v", got)
					if fail.Load() {
						return nil, connect.NewError(connect.CodeInvalidArgument, "invalid automation")
					}
					return response, nil
				}})
				write := c.Create
				if method == "UpdateAutomation" {
					write = c.Update
				}
				result, err := write(t.Context(), input)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, automationID, result.ID)
				assert.Equal(t, "Server response", result.Name)
				assert.Equal(t, input.Disabled, result.Disabled)
				assert.Equal(t, input.Prototype, result.Prototype)
				assert.Equal(t, input.CronTriggers, result.CronTriggers)
				if kind == "storage" {
					require.Len(t, result.StorageEventTriggers, 2)
					assert.Equal(t, retainedTriggerID, result.StorageEventTriggers[0].ID)
					assert.Equal(t, &StorageLocation{ID: locationID}, result.StorageEventTriggers[0].StorageLocation)
					assert.Equal(t, "scenes/**/*.tif", result.StorageEventTriggers[0].GlobPattern)
					assert.Equal(t, &StorageLocation{ID: secondLocationID}, result.StorageEventTriggers[1].StorageLocation)
					assert.Equal(t, "metadata/*.json", result.StorageEventTriggers[1].GlobPattern)
				} else {
					assert.Empty(t, result.StorageEventTriggers)
				}
				assert.Equal(t, id, input.ID)
				assert.Equal(t, "Process arrivals", input.Name)
				fail.Store(true)
				result, err = write(t.Context(), input)
				require.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
			})
		}
	}
}

func TestAutomationClientDelete(t *testing.T) {
	id := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1546")
	for _, cancelJobs := range []bool{false, true} {
		var fail atomic.Bool
		c := newAutomationTestClient(t, "DeleteAutomation", automationRPCService{deleteAutomation: func(_ context.Context, request *workflowsv1.DeleteAutomationRequest) (*emptypb.Empty, error) {
			assert.Equal(t, id, request.GetAutomationId().AsUUID())
			assert.Equal(t, cancelJobs, request.GetCancelJobs())
			if fail.Load() {
				return nil, connect.NewError(connect.CodeNotFound, "missing automation")
			}
			return &emptypb.Empty{}, nil
		}})
		require.NoError(t, c.Delete(t.Context(), id, cancelJobs))
		fail.Store(true)
		err := c.Delete(t.Context(), id, cancelJobs)
		require.Error(t, err)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
	}
}

type automationRPCService struct {
	workflowsv1connect.UnimplementedAutomationServiceHandler

	write            func(context.Context, *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error)
	deleteAutomation func(context.Context, *workflowsv1.DeleteAutomationRequest) (*emptypb.Empty, error)
}

func (s automationRPCService) CreateAutomation(ctx context.Context, request *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error) {
	return s.write(ctx, request)
}

func (s automationRPCService) UpdateAutomation(ctx context.Context, request *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error) {
	return s.write(ctx, request)
}

func (s automationRPCService) DeleteAutomation(ctx context.Context, request *workflowsv1.DeleteAutomationRequest) (*emptypb.Empty, error) {
	return s.deleteAutomation(ctx, request)
}

func newAutomationTestClient(t *testing.T, method string, service automationRPCService) AutomationClient {
	t.Helper()
	rpcServer := connect.NewServer()
	workflowsv1connect.RegisterAutomationServiceHandler(rpcServer, service)
	mux := http.NewServeMux()
	connecthttp.Mount(mux, rpcServer)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/workflows.v1.AutomationService/"+method, r.URL.Path)
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return NewClient(client.New(client.WithURL(server.URL), client.WithHTTPClient(server.Client()), client.WithAPIKey("test-key")), WithDisableTracing()).Automations
}

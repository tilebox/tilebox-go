package workflows

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/workflows/v1/workflowsv1connect"
	"github.com/tilebox/tilebox-go/query"
)

type environmentVariableAPI struct {
	workflowsv1connect.UnimplementedResourceServiceHandler

	get    func(*workflowsv1.GetEnvironmentVariableRequest) (*workflowsv1.EnvironmentVariable, error)
	update func(*workflowsv1.UpdateEnvironmentVariableRequest) (*workflowsv1.EnvironmentVariable, error)
	list   func(*workflowsv1.ListEnvironmentVariablesRequest) (*workflowsv1.ListEnvironmentVariablesResponse, error)
	delete func(*workflowsv1.DeleteEnvironmentVariableRequest) (*workflowsv1.DeleteEnvironmentVariableResponse, error)
}

func (s *environmentVariableAPI) GetEnvironmentVariable(_ context.Context, req *workflowsv1.GetEnvironmentVariableRequest) (*workflowsv1.EnvironmentVariable, error) {
	return s.get(req)
}

func (s *environmentVariableAPI) UpdateEnvironmentVariable(_ context.Context, req *workflowsv1.UpdateEnvironmentVariableRequest) (*workflowsv1.EnvironmentVariable, error) {
	return s.update(req)
}

func (s *environmentVariableAPI) ListEnvironmentVariables(_ context.Context, req *workflowsv1.ListEnvironmentVariablesRequest) (*workflowsv1.ListEnvironmentVariablesResponse, error) {
	return s.list(req)
}

func (s *environmentVariableAPI) DeleteEnvironmentVariable(_ context.Context, req *workflowsv1.DeleteEnvironmentVariableRequest) (*workflowsv1.DeleteEnvironmentVariableResponse, error) {
	return s.delete(req)
}

func newEnvironmentVariableTestClient(t *testing.T, service *environmentVariableAPI) EnvironmentVariableClient {
	t.Helper()
	server := connect.NewServer()
	workflowsv1connect.RegisterResourceServiceHandler(server, service)
	mux := http.NewServeMux()
	connecthttp.Mount(mux, server)
	httpServer := httptest.NewServer(mux)
	t.Cleanup(httpServer.Close)
	rpc := connect.NewClient(connecthttp.NewTransport(httpServer.Client(), httpServer.URL))
	return NewClient(rpc, WithDisableTracing()).EnvironmentVariables
}

func TestEnvironmentVariableCreateJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/workflows.v1.ResourceService/CreateEnvironmentVariable", r.URL.Path)
		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{
			"name": "API_KEY", "value": " literal\n$VALUE ", "isSecret": true,
			"scope": map[string]any{"clusterSlugs": []any{""}}, "description": "Reader-visible description",
		}, body)
		w.Header().Set("Content-Type", "application/json")
		_, err := io.WriteString(w, `{"id":{"uuid":"AAAAAAAAAAAAAAAAAAAAAQ=="},"name":"API_KEY","isSecret":true,"scope":{"clusterSlugs":[""]},"description":"Reader-visible description"}`)
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	rpc := connect.NewClient(connecthttp.NewTransport(server.Client(), server.URL, connecthttp.WithProtoJSON()))
	client := NewClient(rpc, WithDisableTracing()).EnvironmentVariables
	v, err := client.Create(t.Context(), CreateEnvironmentVariableParams{
		Name: "API_KEY", Value: " literal\n$VALUE ", IsSecret: true,
		Scope: ClusterScope{ClusterSlugs: []string{""}}, Description: "Reader-visible description",
	})
	require.NoError(t, err)
	assert.Equal(t, uuid.MustParse("00000000-0000-0000-0000-000000000001"), v.ID)
	assert.Equal(t, "API_KEY", v.Name)
	assert.True(t, v.IsSecret)
	assert.Equal(t, ClusterScope{ClusterSlugs: []string{""}}, v.Scope)
	assert.Equal(t, "Reader-visible description", v.Description)
	assert.Nil(t, v.Value)
	assert.Nil(t, v.SecretPreview)
}

func TestEnvironmentVariableDelete(t *testing.T) {
	id := uuid.New()
	for _, code := range []connect.Code{0, connect.CodePermissionDenied} {
		client := newEnvironmentVariableTestClient(t, &environmentVariableAPI{
			delete: func(req *workflowsv1.DeleteEnvironmentVariableRequest) (*workflowsv1.DeleteEnvironmentVariableResponse, error) {
				assert.Equal(t, id, req.GetId().AsUUID())
				if code != 0 {
					return nil, connect.NewError(code, "access denied")
				}
				return &workflowsv1.DeleteEnvironmentVariableResponse{}, nil
			},
		})
		err := client.Delete(t.Context(), id)
		if code == 0 {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			assert.Equal(t, code, connect.CodeOf(err))
		}
	}
}

func TestEnvironmentVariableOptionalUpdates(t *testing.T) {
	id := uuid.New()
	for _, params := range []UpdateEnvironmentVariableParams{
		{Value: new(" replacement\n$LITERAL ")},
		{Name: new("RENAMED"), IsSecret: new(false), Description: new("")},
		{Scope: &ClusterScope{AllClusters: true}},
		{Scope: &ClusterScope{ClusterSlugs: []string{"b", "a", ""}}},
	} {
		var received *workflowsv1.UpdateEnvironmentVariableRequest
		client := newEnvironmentVariableTestClient(t, &environmentVariableAPI{
			update: func(req *workflowsv1.UpdateEnvironmentVariableRequest) (*workflowsv1.EnvironmentVariable, error) {
				received = req
				return workflowsv1.EnvironmentVariable_builder{Id: tileboxv1.NewUUID(id), Name: "API_KEY", IsSecret: true}.Build(), nil
			},
		})
		result, err := client.Update(t.Context(), id, params)
		require.NoError(t, err)
		assert.Equal(t, id, result.ID)
		assert.Nil(t, result.Value)
		assert.Nil(t, result.SecretPreview)
		require.NotNil(t, received)
		assert.Equal(t, id, received.GetId().AsUUID())
		assert.Equal(t, params.Value != nil, received.HasValue())
		assert.Equal(t, params.Name != nil, received.HasName())
		assert.Equal(t, params.Description != nil, received.HasDescription())
		assert.Equal(t, params.IsSecret != nil, received.HasIsSecret())
		assert.Equal(t, params.Scope != nil, received.HasScope())
		if params.Value != nil {
			assert.Equal(t, *params.Value, received.GetValue())
		}
		if params.Name != nil {
			assert.Equal(t, "RENAMED", received.GetName())
			assert.Empty(t, received.GetDescription())
			assert.False(t, received.GetIsSecret())
		}
		if params.Scope != nil {
			assert.Equal(t, params.Scope.AllClusters, received.GetScope().GetAllClusters())
			assert.Equal(t, params.Scope.ClusterSlugs, received.GetScope().GetClusterSlugs())
		}
	}
}

func TestEnvironmentVariableDisclosure(t *testing.T) {
	id := uuid.New()
	client := newEnvironmentVariableTestClient(t, &environmentVariableAPI{
		get: func(req *workflowsv1.GetEnvironmentVariableRequest) (*workflowsv1.EnvironmentVariable, error) {
			assert.Equal(t, id, req.GetId().AsUUID())
			v := workflowsv1.EnvironmentVariable_builder{Id: tileboxv1.NewUUID(id), Name: "API_KEY", IsSecret: true, Description: "reader description", Scope: workflowsv1.ClusterScope_builder{ClusterSlugs: []string{"a", "b"}}.Build()}.Build()
			if req.GetIncludeSecret() {
				v.SetValue("full-test-value")
			} else {
				v.SetSecretPreview("***")
			}
			return v, nil
		},
	})
	for _, reveal := range []bool{false, true} {
		v, err := client.Get(t.Context(), id, reveal)
		require.NoError(t, err)
		assert.Equal(t, "reader description", v.Description)
		assert.Equal(t, []string{"a", "b"}, v.Scope.ClusterSlugs)
		assert.True(t, v.IsSecret)
		if reveal {
			assert.Equal(t, new("full-test-value"), v.Value)
			assert.Nil(t, v.SecretPreview)
		} else {
			assert.Nil(t, v.Value)
			assert.Equal(t, new("***"), v.SecretPreview)
		}
	}
}

func TestEnvironmentVariableIterationStopsAndReturnsErrors(t *testing.T) {
	id := uuid.New()
	client := newEnvironmentVariableTestClient(t, &environmentVariableAPI{
		list: func(req *workflowsv1.ListEnvironmentVariablesRequest) (*workflowsv1.ListEnvironmentVariablesResponse, error) {
			if req.GetPage().GetStartingAfter() != nil {
				return nil, connect.NewError(connect.CodePermissionDenied, "next page denied")
			}
			return workflowsv1.ListEnvironmentVariablesResponse_builder{
				EnvironmentVariables: []*workflowsv1.EnvironmentVariable{
					workflowsv1.EnvironmentVariable_builder{Id: tileboxv1.NewUUID(id), Name: "FIRST"}.Build(),
				},
				NextPage: tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(id)}.Build(),
			}.Build(), nil
		},
	})
	for _, stopEarly := range []bool{true, false} {
		var ids []uuid.UUID
		var failures []error
		for v, err := range client.List(t.Context(), ListEnvironmentVariablesParams{}) {
			if err != nil {
				assert.Nil(t, v)
				failures = append(failures, err)
				continue
			}
			ids = append(ids, v.ID)
			if stopEarly {
				break
			}
		}
		assert.Equal(t, []uuid.UUID{id}, ids)
		if stopEarly {
			assert.Empty(t, failures)
		} else {
			require.Len(t, failures, 1)
			assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(failures[0]))
		}
	}
}

func TestEnvironmentVariablePagination(t *testing.T) {
	ids := make([]uuid.UUID, 252)
	for i := range ids {
		ids[i] = uuid.New()
	}
	for _, cluster := range []*string{nil, new(""), new("cluster-a")} {
		for _, limit := range []int64{0, 1, 2, 150, 250, 252} {
			client := newEnvironmentVariableTestClient(t, &environmentVariableAPI{
				list: func(req *workflowsv1.ListEnvironmentVariablesRequest) (*workflowsv1.ListEnvironmentVariablesResponse, error) {
					assert.Equal(t, cluster != nil, req.HasClusterSlug())
					if cluster != nil {
						assert.Equal(t, *cluster, req.GetClusterSlug())
					}
					assert.True(t, req.GetIncludeSecrets())
					assert.LessOrEqual(t, req.GetPage().GetLimit(), int64(100))
					index := 0
					if cursor := req.GetPage().GetStartingAfter(); cursor != nil {
						for i, id := range ids {
							if cursor.AsUUID() == id {
								index = i + 1
							}
						}
					}
					pageSize := req.GetPage().GetLimit()
					if limit > 0 {
						assert.Equal(t, min(limit-int64(index-1), 100), pageSize)
					}
					if pageSize == 0 {
						pageSize = 100
					}
					end := min(index+int(pageSize), len(ids))
					variables := make([]*workflowsv1.EnvironmentVariable, 0, end-index)
					for _, id := range ids[index:end] {
						variables = append(variables, workflowsv1.EnvironmentVariable_builder{Id: tileboxv1.NewUUID(id), Name: "TEST", Value: new("test-value")}.Build())
					}
					response := workflowsv1.ListEnvironmentVariablesResponse_builder{EnvironmentVariables: variables}.Build()
					if end < len(ids) {
						response.SetNextPage(tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(ids[end-1])}.Build())
					}
					return response, nil
				},
			})
			params := ListEnvironmentVariablesParams{ClusterSlug: cluster, IncludeSecrets: true, Limit: limit, Cursor: query.NewCursor(ids[0])}
			sequence := client.List(t.Context(), params)
			for range 2 { // Iterating again must start at the original cursor.
				var got []uuid.UUID
				for v, err := range sequence {
					require.NoError(t, err)
					got = append(got, v.ID)
					assert.Equal(t, new("test-value"), v.Value)
				}
				want := ids[1:]
				if limit > 0 && limit < int64(len(want)) {
					want = want[:limit]
				}
				assert.Equal(t, want, got)
			}
		}
	}
}

func TestEnvironmentVariableErrorsPreserveCodeAndOperation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, err := io.WriteString(w, `{"code":"permission_denied","message":"access denied"}`)
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	rpc := connect.NewClient(connecthttp.NewTransport(server.Client(), server.URL))
	client := NewClient(rpc, WithDisableTracing()).EnvironmentVariables
	id := uuid.New()
	for _, tc := range []struct {
		operation string
		call      func() error
	}{
		{"create environment variable", func() error {
			_, err := client.Create(t.Context(), CreateEnvironmentVariableParams{Name: "TOKEN", Value: "private-value"})
			return err
		}},
		{"get environment variable", func() error {
			_, err := client.Get(t.Context(), id, true)
			return err
		}},
		{"list environment variables", func() error {
			_, err := client.ListPage(t.Context(), ListEnvironmentVariablesParams{})
			return err
		}},
		{"update environment variable", func() error {
			_, err := client.Update(t.Context(), id, UpdateEnvironmentVariableParams{Value: new("private-value")})
			return err
		}},
		{"delete environment variable", func() error { return client.Delete(t.Context(), id) }},
	} {
		t.Run(tc.operation, func(t *testing.T) {
			err := tc.call()
			require.EqualError(t, err, "failed to "+tc.operation+": permission_denied: access denied")
			assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
		})
	}
}

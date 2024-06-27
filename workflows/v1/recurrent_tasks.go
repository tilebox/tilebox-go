package workflows

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RecurrentTaskService struct {
	client workflowsv1connect.RecurrentTaskServiceClient
}

func NewRecurrentTaskService(client workflowsv1connect.RecurrentTaskServiceClient) *RecurrentTaskService {
	return &RecurrentTaskService{
		client: client,
	}
}

func (rs *RecurrentTaskService) ListBuckets(ctx context.Context) ([]*workflowsv1.Bucket, error) {
	response, err := rs.client.ListBuckets(ctx, connect.NewRequest(&emptypb.Empty{}))

	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}

	return response.Msg.GetBuckets(), nil
}

func (rs *RecurrentTaskService) GetBucket(ctx context.Context, bucketID *workflowsv1.UUID) (*workflowsv1.Bucket, error) {
	response, err := rs.client.GetBucket(ctx, connect.NewRequest(bucketID))

	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	return response.Msg, nil
}

func (rs *RecurrentTaskService) CreateBucket(ctx context.Context, name string, bucketType workflowsv1.BucketType) (*workflowsv1.Bucket, error) {
	req := connect.NewRequest(&workflowsv1.Bucket{
		Name: name,
		Type: bucketType,
	})
	response, err := rs.client.CreateBucket(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return response.Msg, nil
}

func (rs *RecurrentTaskService) DeleteBucket(ctx context.Context, bucketID *workflowsv1.UUID) error {
	_, err := rs.client.DeleteBucket(ctx, connect.NewRequest(bucketID))

	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}

	return nil
}

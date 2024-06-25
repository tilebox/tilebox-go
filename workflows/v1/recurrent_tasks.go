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
		return nil, fmt.Errorf("failed to submit job: %w", err)
	}

	return response.Msg.GetBuckets(), nil
}

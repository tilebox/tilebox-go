package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"github.com/google/uuid"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
)

func protoToUUID(id *workflowsv1.UUID) (uuid.UUID, error) {
	if id == nil || len(id.GetUuid()) == 0 {
		return uuid.Nil, nil
	}

	bytes, err := uuid.FromBytes(id.GetUuid())
	if err != nil {
		return uuid.Nil, err
	}

	return bytes, nil
}

func uuidToProtobuf(id uuid.UUID) *workflowsv1.UUID {
	if id == uuid.Nil {
		return nil
	}

	return &workflowsv1.UUID{Uuid: id[:]}
}

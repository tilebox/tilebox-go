// Code is NOT generated.

package tileboxv1

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// NewUUID constructs a new ID from the provided uuid.UUID.
func NewUUID(id uuid.UUID) *ID {
	if id == uuid.Nil {
		return nil
	}

	return ID_builder{Uuid: id[:]}.Build()
}

func NewUUIDSlice(ids []uuid.UUID) []*ID {
	pbIDs := make([]*ID, 0, len(ids))
	for _, id := range ids {
		pbIDs = append(pbIDs, NewUUID(id))
	}
	return pbIDs
}

// AsUUID converts x to a uuid.UUID.
func (x *ID) AsUUID() uuid.UUID {
	if x == nil || len(x.GetUuid()) == 0 {
		return uuid.Nil
	}

	id, err := uuid.FromBytes(x.GetUuid())
	if err != nil {
		return uuid.Nil
	}
	return id
}

// IsValid reports whether the uuid is valid.
// It is equivalent to CheckValid == nil.
func (x *ID) IsValid() bool {
	return x.CheckValid() == nil
}

// CheckValid returns an error if the uuid is invalid.
// An error is reported for a nil UUID.
func (x *ID) CheckValid() error {
	if x == nil {
		return protoimpl.X.NewError("invalid UUID: nil")
	}

	_, err := uuid.FromBytes(x.GetUuid())
	if err != nil {
		return protoimpl.X.NewError("invalid UUID: %v", err)
	}
	return nil
}

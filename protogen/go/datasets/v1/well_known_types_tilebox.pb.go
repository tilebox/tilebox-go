// Code is NOT generated.

package datasetsv1

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// NewUUID constructs a new UUID from the provided uuid.UUID.
func NewUUID(id uuid.UUID) *UUID {
	return &UUID{Uuid: id[:]}
}

// AsUUID converts x to a uuid.UUID.
func (x *UUID) AsUUID() uuid.UUID {
	id, err := uuid.FromBytes(x.GetUuid())
	if err != nil {
		return uuid.Nil
	}
	return id
}

// IsValid reports whether the timestamp is valid.
// It is equivalent to CheckValid == nil.
func (x *UUID) IsValid() bool {
	return x.CheckValid() == nil
}

// CheckValid returns an error if the uuid is invalid.
// An error is reported for a nil UUID.
func (x *UUID) CheckValid() error {
	if x == nil {
		return protoimpl.X.NewError("invalid UUID: nil")
	}

	_, err := uuid.FromBytes(x.GetUuid())
	if err != nil {
		return protoimpl.X.NewError("invalid UUID: %v", err)
	}
	return nil
}

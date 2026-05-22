package query

import (
	"fmt"

	"github.com/google/uuid"
)

// Cursor identifies where to continue a paginated query.
//
// Cursors are returned as NextCursor values from page query methods and should only be reused with the same endpoint,
// filters and sort direction that produced them.
type Cursor struct {
	startingAfter uuid.UUID
}

// NewCursor creates a cursor that starts after the entry with the given ID.
func NewCursor(startingAfter uuid.UUID) *Cursor {
	return &Cursor{startingAfter: startingAfter}
}

// ParseCursor parses a cursor string returned by Cursor.String.
func ParseCursor(value string) (*Cursor, error) {
	startingAfter, err := uuid.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}
	return NewCursor(startingAfter), nil
}

// String returns the cursor as a string that can be passed to ParseCursor.
func (c *Cursor) String() string {
	if c == nil {
		return ""
	}
	return c.startingAfter.String()
}

// StartingAfter returns the ID after which the next page should start.
func (c *Cursor) StartingAfter() uuid.UUID {
	if c == nil {
		return uuid.Nil
	}
	return c.startingAfter
}

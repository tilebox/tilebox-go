package workflow

import (
	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/query"
)

// ListOptions contains workflow list pagination options.
type ListOptions struct {
	Limit  int64
	Cursor *Cursor
}

// ListOption configures a workflow list or page request.
type ListOption interface {
	ApplyListOption(*ListOptions)
}

// Cursor identifies where to continue a paginated workflow list.
type Cursor = query.Cursor

// NewCursor creates a cursor that starts after the workflow with the given ID.
func NewCursor(startingAfter uuid.UUID) *Cursor {
	return query.NewCursor(startingAfter)
}

// ParseCursor parses a cursor string returned by Cursor.String.
func ParseCursor(value string) (*Cursor, error) {
	return query.ParseCursor(value)
}

type listOptionFunc func(*ListOptions)

func (f listOptionFunc) ApplyListOption(options *ListOptions) {
	f(options)
}

// WithCursor starts the list after a cursor returned by the same list endpoint.
func WithCursor(cursor *Cursor) ListOption {
	return listOptionFunc(func(options *ListOptions) { options.Cursor = cursor })
}

// WithLimit limits the total workflows returned by an iterator, or the size of a single page.
// A non-positive limit uses the server default for pages and no total limit for iterators.
func WithLimit(limit int64) ListOption {
	return listOptionFunc(func(options *ListOptions) { options.Limit = limit })
}

// NewListOptions applies workflow list options.
func NewListOptions(options ...ListOption) ListOptions {
	var applied ListOptions
	for _, option := range options {
		option.ApplyListOption(&applied)
	}
	return applied
}

// UndeployOptions contains workflow undeployment options.
type UndeployOptions struct {
	ReleaseID *uuid.UUID
}

// UndeployOption configures workflow undeployment requests.
type UndeployOption interface {
	ApplyUndeployOption(*UndeployOptions)
}

type releaseIDOption struct {
	releaseID uuid.UUID
}

func (o releaseIDOption) ApplyUndeployOption(options *UndeployOptions) {
	options.ReleaseID = &o.releaseID
}

// WithReleaseID restricts undeployment to the given release instead of all releases.
func WithReleaseID(releaseID uuid.UUID) UndeployOption {
	return releaseIDOption{releaseID: releaseID}
}

// NewUndeployOptions applies undeployment options.
func NewUndeployOptions(options ...UndeployOption) UndeployOptions {
	var applied UndeployOptions
	for _, option := range options {
		option.ApplyUndeployOption(&applied)
	}
	return applied
}

// UpdateOptions contains workflow update options.
type UpdateOptions struct {
	Name        *string
	Description *string
}

// UpdateOption configures workflow update requests.
type UpdateOption interface {
	ApplyUpdateOption(*UpdateOptions)
}

type nameOption struct {
	name string
}

func (o nameOption) ApplyUpdateOption(options *UpdateOptions) {
	options.Name = &o.name
}

type descriptionOption struct {
	description string
}

func (o descriptionOption) ApplyUpdateOption(options *UpdateOptions) {
	options.Description = &o.description
}

// WithName sets the workflow name.
func WithName(name string) UpdateOption {
	return nameOption{name: name}
}

// WithDescription sets the workflow description.
func WithDescription(description string) UpdateOption {
	return descriptionOption{description: description}
}

// NewUpdateOptions applies update options.
func NewUpdateOptions(options ...UpdateOption) UpdateOptions {
	var applied UpdateOptions
	for _, option := range options {
		option.ApplyUpdateOption(&applied)
	}
	return applied
}

// Package storagelocation provides options for storage subscription event queries.
package storagelocation

import "github.com/tilebox/tilebox-go/query"

// ListOptions contains storage subscription event pagination options.
type ListOptions struct {
	Limit  int64
	Cursor *query.Cursor
}

// ListOption configures an event list or page request.
type ListOption func(*ListOptions)

// WithCursor continues from a cursor returned by the same storage location's event list.
func WithCursor(cursor *query.Cursor) ListOption {
	return func(options *ListOptions) { options.Cursor = cursor }
}

// WithLimit limits the total events returned by an iterator, or the size of a single page.
// Pages are capped at 100. A non-positive limit uses the server default for pages
// and no total limit for iterators.
func WithLimit(limit int64) ListOption {
	return func(options *ListOptions) { options.Limit = limit }
}

// NewListOptions applies event list options.
func NewListOptions(options ...ListOption) ListOptions {
	var applied ListOptions
	for _, option := range options {
		option(&applied)
	}
	return applied
}

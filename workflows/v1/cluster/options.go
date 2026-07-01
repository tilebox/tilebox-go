package cluster // import "github.com/tilebox/tilebox-go/workflows/v1/cluster"

// CreateOptions contains cluster create options.
type CreateOptions struct {
	Description string
	Slug        string
}

// UpdateOptions contains cluster update options.
type UpdateOptions struct {
	Name        *string
	Description *string
}

// CreateOption configures cluster create requests.
type CreateOption interface {
	ApplyCreateOption(*CreateOptions)
}

// UpdateOption configures cluster update requests.
type UpdateOption interface {
	ApplyUpdateOption(*UpdateOptions)
}

// DescriptionOption configures cluster descriptions on create and update requests.
type DescriptionOption interface {
	CreateOption
	UpdateOption
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

func (o descriptionOption) ApplyCreateOption(options *CreateOptions) {
	options.Description = o.description
}

func (o descriptionOption) ApplyUpdateOption(options *UpdateOptions) {
	options.Description = &o.description
}

type slugOption struct {
	slug string
}

func (o slugOption) ApplyCreateOption(options *CreateOptions) {
	options.Slug = o.slug
}

// WithName sets the cluster name for update requests.
func WithName(name string) UpdateOption {
	return nameOption{name: name}
}

// WithDescription sets the cluster description.
func WithDescription(description string) DescriptionOption {
	return descriptionOption{description: description}
}

// WithSlug sets the optional cluster slug for create requests.
func WithSlug(slug string) CreateOption {
	return slugOption{slug: slug}
}

// NewCreateOptions applies create options.
func NewCreateOptions(options ...CreateOption) CreateOptions {
	var applied CreateOptions
	for _, option := range options {
		option.ApplyCreateOption(&applied)
	}
	return applied
}

// NewUpdateOptions applies update options.
func NewUpdateOptions(options ...UpdateOption) UpdateOptions {
	var applied UpdateOptions
	for _, option := range options {
		option.ApplyUpdateOption(&applied)
	}
	return applied
}

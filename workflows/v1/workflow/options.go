package workflow

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

package filter

import (
	"strconv"
)

// FilterBuilder provides a fluent interface for building query filters
// Similar to Laravel's query builder pattern
type FilterBuilder struct {
	filters map[string]any
}

// NewFilterBuilder creates a new filter builder
func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		filters: make(map[string]any),
	}
}

// String adds a string filter if the value is not empty
func (b *FilterBuilder) String(key, value string) *FilterBuilder {
	if value != "" {
		b.filters[key] = value
	}
	return b
}

// Uint adds a uint filter from a string value
// Returns the builder for chaining even if parsing fails
func (b *FilterBuilder) Uint(key, value string) *FilterBuilder {
	if value != "" {
		if parsed, err := strconv.ParseUint(value, 10, 32); err == nil {
			b.filters[key] = uint(parsed)
		}
	}
	return b
}

// UintValue adds a uint filter directly
func (b *FilterBuilder) UintValue(key string, value uint) *FilterBuilder {
	if value > 0 {
		b.filters[key] = value
	}
	return b
}

// Int adds an int filter from a string value
func (b *FilterBuilder) Int(key, value string) *FilterBuilder {
	if value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			b.filters[key] = parsed
		}
	}
	return b
}

// IntValue adds an int filter directly
func (b *FilterBuilder) IntValue(key string, value int) *FilterBuilder {
	b.filters[key] = value
	return b
}

// Bool adds a bool filter from a string value
func (b *FilterBuilder) Bool(key, value string) *FilterBuilder {
	if value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			b.filters[key] = parsed
		}
	}
	return b
}

// BoolValue adds a bool filter directly
func (b *FilterBuilder) BoolValue(key string, value bool) *FilterBuilder {
	b.filters[key] = value
	return b
}

// Where adds a filter with any value type
func (b *FilterBuilder) Where(key string, value any) *FilterBuilder {
	if value != nil {
		b.filters[key] = value
	}
	return b
}

// WhereIn adds an IN filter (for slice values)
func (b *FilterBuilder) WhereIn(key string, values []any) *FilterBuilder {
	if len(values) > 0 {
		b.filters[key] = values
	}
	return b
}

// Build returns the built filters map
func (b *FilterBuilder) Build() map[string]any {
	return b.filters
}

// Get returns a specific filter value
func (b *FilterBuilder) Get(key string) (any, bool) {
	val, exists := b.filters[key]
	return val, exists
}

// Has checks if a filter exists
func (b *FilterBuilder) Has(key string) bool {
	_, exists := b.filters[key]
	return exists
}

// Count returns the number of filters
func (b *FilterBuilder) Count() int {
	return len(b.filters)
}

// Clear removes all filters
func (b *FilterBuilder) Clear() *FilterBuilder {
	b.filters = make(map[string]any)
	return b
}

// Remove removes a specific filter
func (b *FilterBuilder) Remove(key string) *FilterBuilder {
	delete(b.filters, key)
	return b
}

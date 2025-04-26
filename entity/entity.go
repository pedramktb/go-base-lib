package entity

// Entity is the base interface for all domain models.
// We use method implementation instead of struct tags for performance reasons.
type Entity interface {
	// New is a constructor method for creating new non-nil entities.
	New() Entity
	// Fields should be implemented to provide a map of all fields to their pointers.
	Fields() (fields map[string]any)
}

// UpdateEntity is the base interface for all domain update models.
// Optional values in update models are represented by the `Optional` type.
type UpdateEntity interface {
	// Fields should be implemented to provide the optional fields values and whether they are set or not.
	Field(field string) (value any, set bool)
}

type e[t m] interface {
	New() t
	m
}

type m interface {
	// Fields should be implemented to provide a map of all fields to their pointers.
	Fields() (fields map[string]any)
}

type f interface {
	e[f]
}

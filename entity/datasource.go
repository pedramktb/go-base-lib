package entity

import (
	"context"

	"github.com/pedramktb/go-base-lib/types"
)

// Getter is the interface for entities that can be retrieved by their unique identifier.
type Getter[E Entity] interface {
	Get(ctx context.Context, id E) (E, error)
}

// Lister is the interface for listing entities.
// Intended for those a limited count with no pagination or filters.
type Lister[E Entity] interface {
	List(ctx context.Context) ([]E, error)
}

// Querier is the interface for querying entities with potential filters, searching, and sorts.
// Intended for those with a large count and pagination.
type Querier[E Entity] interface {
	Query(
		ctx context.Context,
		filters Expression[E],
		search types.Nillable[string],
		sorts Sorts[E],
		limit PaginationLimit,
	) (Paginated[E], error)
}

// Creator is the interface for creating entities.
type Creator[E Entity] interface {
	Create(ctx context.Context, items ...E) error
}

// Updater is the interface for updating an entity.
type Updater[U UpdateEntity] interface {
	Update(ctx context.Context, update U) error
}

// Deleter is the interface for deleting entities.
type Deleter[E Entity] interface {
	Delete(ctx context.Context, id E) error
}

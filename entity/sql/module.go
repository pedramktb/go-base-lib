package entitysql

import (
	"context"
	"database/sql"

	"github.com/pedramktb/go-base-lib/entity"
	"github.com/pedramktb/go-base-lib/env"
	"github.com/pressly/goose/v3/database"
)

// Returns sql implementation of given Entity and table Name
func Module[E SQLEntity, U entity.UpdateEntity](ctx context.Context, db *sql.DB, dbDialect database.Dialect, table string) (
	entity.Getter[E],
	entity.Lister[E],
	entity.Querier[E],
	entity.Creator[E],
	entity.Updater[U],
	entity.Deleter[E],
	error,
) {
	module, err := NewEntitySQL[E, U](ctx, db, dbDialect, table, "migrations/"+env.GetEnvironment().String())
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	return module, module, module, module, module, module, nil
}

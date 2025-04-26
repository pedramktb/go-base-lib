package entitysql

import (
	"context"
	"database/sql"
	"errors"

	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pedramktb/go-base-lib/entity"
	"github.com/pedramktb/go-base-lib/taggederror"
	"github.com/pedramktb/go-base-lib/types"
	"github.com/pressly/goose/v3/database"
)

type Column struct {
	Type       string
	PrimaryKey bool
	Nullable   bool
}

// SQLEntity adds PKs for fetching primary keys and Row for scanning and inserting rows.
type SQLEntity interface {
	entity.Entity
	Row() []any
	Columns() (columns map[string]Column)
}

type entitySQL[E SQLEntity, U entity.UpdateEntity] struct {
	db      *sql.DB
	table   string
	columns []string
}

func NewEntitySQL[E SQLEntity, U entity.UpdateEntity](
	ctx context.Context,
	db *sql.DB,
	dialect database.Dialect,
	migrationDir,
	table string,
) (*entitySQL[E, U], error) {
	if err := migrate[E](ctx, db, dialect, table, migrationDir); err != nil {
		return nil, err
	}
	columns := make([]string, 0, len((*new(E)).Columns()))
	for k := range (*new(E)).Columns() {
		columns = append(columns, k)
	}
	return &entitySQL[E, U]{
		db,
		table,
		columns,
	}, nil
}

func (db *entitySQL[E, U]) locateWithID(id E) squirrel.Sqlizer {
	var locate squirrel.Sqlizer
	for field, column := range id.Columns() {
		if !column.PrimaryKey {
			continue
		}
		if locate == nil {
			locate = squirrel.Eq{field: id.Fields()[field]}
		} else {
			locate = squirrel.And{locate, squirrel.Eq{field: id.Fields()[field]}}
		}
	}
	return locate
}

func (db *entitySQL[E, U]) locateUpdateWithID(id U) squirrel.Sqlizer {
	var locate squirrel.Sqlizer
	for field, column := range (*new(E)).Columns() {
		if !column.PrimaryKey {
			continue
		}
		value, _ := id.Field(field)
		if locate == nil {
			locate = squirrel.Eq{field: value}
		} else {
			locate = squirrel.And{locate, squirrel.Eq{field: value}}
		}
	}
	return locate
}

func (db *entitySQL[E, U]) Get(ctx context.Context, id E) (E, error) {
	var result = (*new(E)).New().(E)

	query, args, err := squirrel.
		Select(db.columns...).
		PlaceholderFormat(squirrel.Dollar).
		From(db.table).
		Where(db.locateWithID(id)).
		ToSql()
	if err != nil {
		return result, entity.ErrDBUnhandled.Wrap(err)
	}

	err = db.db.QueryRowContext(ctx, query, args...).Scan(result.Row()...)
	if errors.Is(err, sql.ErrNoRows) {
		return result, taggederror.ErrNotFound
	} else if err != nil {
		return result, entity.ErrDBUnhandled.Wrap(err)
	}

	return result, nil
}

func (db *entitySQL[E, U]) List(ctx context.Context) ([]E, error) {
	query, args, err := squirrel.
		Select(db.columns...).
		PlaceholderFormat(squirrel.Dollar).
		From(db.table).
		ToSql()
	if err != nil {
		return nil, entity.ErrDBUnhandled.Wrap(err)
	}

	rows, err := db.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, entity.ErrDBUnhandled.Wrap(err)
	}
	defer rows.Close()

	var result []E
	for rows.Next() {
		var e E
		err := rows.Scan(e.Row()...)
		if err != nil {
			return nil, entity.ErrDBUnhandled.Wrap(err)
		}
		result = append(result, e)
	}
	if err := rows.Err(); err != nil {
		return nil, entity.ErrDBUnhandled.Wrap(err)
	}
	return result, nil
}

func (db *entitySQL[E, U]) Query(
	ctx context.Context,
	filters entity.Expression[E],
	search types.Nillable[string],
	sorts entity.Sorts[E],
	limit entity.PaginationLimit,
) (entity.Paginated[E], error) {
	builder := squirrel.
		Select(db.columns...).
		PlaceholderFormat(squirrel.Dollar).
		From(db.table).Where(filters.ToSQLQuery())

	// Apply search to the base query
	if search.NotNil {
		search.Val = strings.ReplaceAll(search.Val, " ", " & ")
		builder = builder.Where("tsv @@ to_tsquery(?)", search.Val)
	}

	// Create a count query based on the base query
	countQuery, countArgs, err := squirrel.
		Select("COUNT(*)").
		PlaceholderFormat(squirrel.Dollar).
		FromSelect(builder, "entity_count").
		ToSql()
	if err != nil {
		return entity.Paginated[E]{}, entity.ErrDBUnhandled.Wrap(err)
	}

	var total uint64
	if err := db.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return entity.Paginated[E]{}, entity.ErrDBUnhandled.Wrap(err)
	}

	builder = sorts.ToSQLQuery(builder).Limit(uint64(limit + 1))

	query, args, err := builder.ToSql()
	if err != nil {
		return entity.Paginated[E]{}, entity.ErrDBUnhandled.Wrap(err)
	}

	rows, err := db.db.QueryContext(ctx, query, args...)
	if err != nil {
		return entity.Paginated[E]{}, entity.ErrDBUnhandled.Wrap(err)
	}
	defer rows.Close()

	result := make([]E, 0, limit+1)
	for rows.Next() {
		var e E
		err := rows.Scan(e.Row()...)
		if err != nil {
			return entity.Paginated[E]{}, entity.ErrDBUnhandled.Wrap(err)
		}
		result = append(result, e)
	}
	if err := rows.Err(); err != nil {
		return entity.Paginated[E]{}, entity.ErrDBUnhandled.Wrap(err)
	}

	var next types.Nillable[E]
	if len(result) > int(limit) {
		next = types.NewNillable(&result[limit-1])
		result = result[:limit-1]
	}

	return entity.Paginated[E]{
		Items: result,
		Meta: entity.PaginationMeta{
			Total: total,
			Next:  sorts.NextCursor(next),
		},
	}, nil
}

func (db *entitySQL[E, U]) Create(ctx context.Context, items ...E) error {
	preQuery := squirrel.
		Insert(db.table).
		PlaceholderFormat(squirrel.Dollar).
		Columns(db.columns...)

	for i := range items {
		preQuery = preQuery.Values(items[i].Row()...)
	}

	insert, args, err := preQuery.ToSql()
	if err != nil {
		return entity.ErrDBUnhandled.Wrap(err)
	}

	_, err = db.db.ExecContext(ctx, insert, args...)
	if err != nil {
		return entity.ErrDBUnhandled.Wrap(err)
	}

	return nil
}

func (db *entitySQL[E, U]) Update(ctx context.Context, update U) error {
	builder := squirrel.
		Update(db.table).
		PlaceholderFormat(squirrel.Dollar).
		Where(db.locateUpdateWithID(update))

	for _, c := range db.columns {
		if val, set := update.Field(c); set {
			builder = builder.Set(c, val)
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return entity.ErrDBUnhandled.Wrap(err)
	}

	_, err = db.db.ExecContext(ctx, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return taggederror.ErrNotFound
	} else if err != nil {
		return entity.ErrDBUnhandled.Wrap(err)
	}

	return nil
}

func (db *entitySQL[E, U]) Delete(ctx context.Context, id E) error {
	query, args, err := squirrel.
		Delete(db.table).
		PlaceholderFormat(squirrel.Dollar).
		Where(db.locateWithID(id)).
		ToSql()
	if err != nil {
		return entity.ErrDBUnhandled.Wrap(err)
	}

	if tag, err := db.db.ExecContext(ctx, query, args...); err == nil {
		if rows, err := tag.RowsAffected(); err == nil && rows == 0 {
			return taggederror.ErrNotFound.Wrap(err)
		}
	} else {
		return entity.ErrDBUnhandled.Wrap(err)
	}

	return nil
}

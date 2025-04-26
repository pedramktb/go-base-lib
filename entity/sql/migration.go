package entitysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
)

func migrate[E SQLEntity](ctx context.Context, db *sql.DB, dialect database.Dialect, table, dir string) error {
	cols := (*new(E)).Columns()

	var up, down string
	oldCols, err := fetchColumns(ctx, db, table)
	if err != nil {
		return err
	}

	if len(oldCols) == 0 {
		up, down = initialMigration(cols, table)
	} else {
		up, down, err = generateDiffMigration(
			oldCols,
			cols,
			table,
		)
		if err != nil {
			return err
		}
	}

	err = writeMigration(dir, table, up, down)
	if err != nil {
		return err
	}

	err = goose.SetDialect(string(dialect))
	if err != nil {
		return err
	}

	return goose.Up(db, dir)
}

func writeMigration(dir, table, up, down string) error {
	filename := fmt.Sprintf("%s/%s_%s.sql", dir, time.Now().UTC().Format("20060102150405"), table)
	content := fmt.Sprintf(
		`-- +goose Up
-- +goose StatementBegin
%s
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
%s
-- +goose StatementEnd
`, up, down)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(content), os.ModePerm)
}

func initialMigration(cols map[string]Column, tableName string) (upSQL, downSQL string) {
	var colDefs []string
	for name, col := range cols {
		def := fmt.Sprintf("%s %s", name, col.Type)
		if col.PrimaryKey {
			def += " PRIMARY KEY"
		}
		if !col.Nullable {
			def += " NOT NULL"
		}
		colDefs = append(colDefs, def)
	}

	upSQL = fmt.Sprintf("CREATE TABLE %s (\n  %s\n);", tableName, strings.Join(colDefs, ",\n  "))
	downSQL = fmt.Sprintf("DROP TABLE %s;", tableName)
	return upSQL, downSQL
}

func generateDiffMigration(oldCols, newCols map[string]Column, tableName string) (upSQL, downSQL string, err error) {
	up := make([]string, 2)
	down := make([]string, 2)
	changed := make([]bool, 2)

	for name, col := range newCols {
		if _, exists := oldCols[name]; !exists {
			if !changed[0] {
				up[0] = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, name, col.Type)
				down[0] = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, name)
				changed[0] = true
			} else {
				up[0] += fmt.Sprintf(", %s %s", name, col.Type)
				down[0] += fmt.Sprintf(", DROP COLUMN %s", name)
			}
		} else if oldCols[name].Type != col.Type {
			return "", "", errors.New("changing column types in migrations is not supported")
		}
	}

	for name := range oldCols {
		if _, exists := newCols[name]; !exists {
			if !changed[1] {
				up[1] = fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", tableName, name)
				down[1] = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` `%s`", tableName, name, oldCols[name].Type)
				changed[1] = true
			} else {
				up[1] += fmt.Sprintf(", DROP COLUMN `%s`", name)
				down[1] += fmt.Sprintf(", ADD COLUMN `%s` `%s`", name, oldCols[name].Type)
			}
		}
	}

	if changed[0] {
		up[0] += ";\n"
		down[0] += ";\n"
	}
	if changed[1] {
		up[1] += ";"
		down[1] += ";"
	}

	return strings.Join(up, ""), strings.Join(down, ""), nil
}

func fetchColumns(ctx context.Context, db *sql.DB, table string) (map[string]Column, error) {
	query, args, err := squirrel.
		Select("COLUMN_NAME").
		PlaceholderFormat(squirrel.Dollar).
		From("information_schema.columns").
		Where(squirrel.Eq{"table_name": table}).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make(map[string]Column)
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt sql.NullString
		err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		if err != nil {
			return nil, err
		}
		cols[name] = Column{Type: ctype, PrimaryKey: pk > 0, Nullable: notnull == 0}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return cols, nil
}

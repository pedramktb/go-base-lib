package entitysql

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/pedramktb/go-base-lib/entity"
	postgrestestcontainer "github.com/pedramktb/go-base-lib/postgres/testcontainer"
	"github.com/pedramktb/go-base-lib/taggederror"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
)

var ip, port string

func TestMain(m *testing.M) {
	ctx := context.Background()
	var postgresContainer testcontainers.Container
	postgresContainer, ip, port = postgrestestcontainer.NewTestContainer(ctx)
	defer func() {
		_ = postgresContainer.Terminate(ctx)
	}()
	m.Run()
}

type testEntity struct {
	ID   int64
	Name string
}

func (*testEntity) New() entity.Entity {
	return new(testEntity)
}

func (e *testEntity) Fields() map[string]any {
	if e == nil {
		e = new(testEntity)
	}
	return map[string]any{
		"id":   &e.ID,
		"name": &e.Name,
	}
}

func (*testEntity) Columns() map[string]Column {
	return map[string]Column{
		"id":   {Type: "BIGINT", PrimaryKey: true},
		"name": {Type: "TEXT"},
	}
}

func (e *testEntity) Row() []any {
	if e == nil {
		e = new(testEntity)
	}
	return []any{&e.ID, &e.Name}
}

type testEntityPatch struct{}

func (testEntityPatch) Field(string) (any, bool) {
	return nil, false
}

func setupDB(t *testing.T, dbName string) (*entitySQL[*testEntity, *testEntityPatch], func()) {
	t.Helper()
	ctx := context.Background()
	db := postgrestestcontainer.CreateTestDB(ctx, ip, port, dbName)
	// apply migrations and seeds
	tearDown := func() {
		postgrestestcontainer.DropTestDB(ctx, db, ip, port, dbName)
		_ = os.RemoveAll("./migrations/test")
	}
	edb, err := NewEntitySQL[*testEntity, *testEntityPatch](
		ctx,
		db,
		goose.DialectPostgres,
		"./migrations/test",
		"test_entity",
	)
	if err != nil {
		t.Fatal(err)
	}
	return edb, tearDown
}

func Test_Get(t *testing.T) {
	edb, tearDown := setupDB(t, "test_entity_get")
	defer tearDown()

	ent1 := &testEntity{ID: 123, Name: "test"}
	err := edb.Create(context.Background(), ent1)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	tests := []struct {
		name    string
		id      *testEntity
		want    *testEntity
		wantErr error
	}{
		{
			name:    "success",
			id:      &testEntity{ID: ent1.ID},
			want:    ent1,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      &testEntity{ID: 3123},
			want:    &testEntity{},
			wantErr: taggederror.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := edb.Get(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if *got != (*tt.want) {
				t.Errorf("Get() got = %v, want %v", *got, *tt.want)
			}
		})
	}
}

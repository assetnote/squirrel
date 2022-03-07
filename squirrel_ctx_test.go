// +build go1.8

package squirrel

import (
	"context"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func (s *DBStub) PrepareContext(ctx context.Context, query string) (*pgconn.StatementDescription, error) {
	s.LastPrepareSql = query
	s.PrepareCount++
	return nil, nil
}

func (s *DBStub) ExecContext(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	s.LastExecSql = query
	s.LastExecArgs = args
	return nil, nil
}

func (s *DBStub) QueryContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	s.LastQuerySql = query
	s.LastQueryArgs = args
	return nil, nil
}

func (s *DBStub) QueryRowContext(ctx context.Context, query string, args ...interface{}) pgx.Row {
	s.LastQueryRowSql = query
	s.LastQueryRowArgs = args
	return &Row{Row: &RowStub{}}
}

var ctx = context.Background()

func TestExecContextWith(t *testing.T) {
	db := &DBStub{}
	ExecContextWith(ctx, db, sqlizer)
	assert.Equal(t, sqlStr, db.LastExecSql)
}

func TestQueryContextWith(t *testing.T) {
	db := &DBStub{}
	QueryContextWith(ctx, db, sqlizer)
	assert.Equal(t, sqlStr, db.LastQuerySql)
}

func TestQueryRowContextWith(t *testing.T) {
	db := &DBStub{}
	QueryRowContextWith(ctx, db, sqlizer)
	assert.Equal(t, sqlStr, db.LastQueryRowSql)
}

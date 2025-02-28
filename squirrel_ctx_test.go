// +build go1.8

package squirrel

import (
	"context"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
)

func (s *DBStub) PrepareContext(ctx context.Context, query string) (*pgconn.StatementDescription, error) {
	s.LastPrepareSql = query
	s.PrepareCount++
	return nil, nil
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

// +build go1.8

package squirrel

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// NoContextSupport is returned if a db doesn't support Context.
var NoContextSupport = errors.New("DB does not support Context")

// WrapStdSqlCtx wraps a type implementing the standard SQL interface plus the context
// versions of the methods with methods that squirrel expects.
func WrapStdSqlCtx(stdSqlCtx StdSql) Runner {
	return &stdsqlCtxRunner{stdSqlCtx}
}

type stdsqlCtxRunner struct {
	StdSql
}

// ExecContextWith ExecContexts the SQL returned by s with db.
func ExecContextWith(ctx context.Context, db Execer, s Sqlizer) (res pgconn.CommandTag, err error) {
	query, args, err := s.ToSql()
	if err != nil {
		return
	}
	return db.Exec(ctx, query, args...)
}

// QueryContextWith QueryContexts the SQL returned by s with db.
func QueryContextWith(ctx context.Context, db Queryer, s Sqlizer) (rows pgx.Rows, err error) {
	query, args, err := s.ToSql()
	if err != nil {
		return
	}
	return db.Query(ctx, query, args...)
}

// QueryRowContextWith QueryRowContexts the SQL returned by s with db.
func QueryRowContextWith(ctx context.Context, db QueryRower, s Sqlizer) pgx.Row {
	query, args, err := s.ToSql()
	return &Row{Row: db.QueryRow(ctx, query, args...), err: err}
}

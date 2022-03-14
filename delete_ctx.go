// +build go1.8

package squirrel

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/lann/builder"
)

func (d *deleteData) ExecContext(ctx context.Context) (pgconn.CommandTag, error) {
	if d.RunWith == nil {
		return nil, RunnerNotSet
	}
	ctxRunner, ok := d.RunWith.(Execer)
	if !ok {
		return nil, NoContextSupport
	}
	return ExecContextWith(ctx, ctxRunner, d)
}

func (d *deleteData) QueryContext(ctx context.Context) (pgx.Rows, error) {
	if d.RunWith == nil {
		return nil, RunnerNotSet
	}
	ctxRunner, ok := d.RunWith.(Queryer)
	if !ok {
		return nil, NoContextSupport
	}
	return QueryContextWith(ctx, ctxRunner, d)
}

func (d *deleteData) QueryRowContext(ctx context.Context) pgx.Row {
	if d.RunWith == nil {
		return &Row{err: RunnerNotSet}
	}
	queryRower, ok := d.RunWith.(QueryRower)
	if !ok {
		if _, ok := d.RunWith.(Queryer); !ok {
			return &Row{err: RunnerNotQueryRunner}
		}
		return &Row{err: NoContextSupport}
	}
	return QueryRowContextWith(ctx, queryRower, d)
}

// ExecContext builds and ExecContexts the query with the Runner set by RunWith.
func (b DeleteBuilder) ExecContext(ctx context.Context) (pgconn.CommandTag, error) {
	data := builder.GetStruct(b).(deleteData)
	return data.ExecContext(ctx)
}

// QueryContext builds and QueryContexts the query with the Runner set by RunWith.
func (b DeleteBuilder) QueryContext(ctx context.Context) (pgx.Rows, error) {
	data := builder.GetStruct(b).(deleteData)
	return data.QueryContext(ctx)
}

// QueryRowContext builds and QueryRowContexts the query with the Runner set by RunWith.
func (b DeleteBuilder) QueryRowContext(ctx context.Context) pgx.Row {
	data := builder.GetStruct(b).(deleteData)
	return data.QueryRowContext(ctx)
}

// ScanContext is a shortcut for QueryRowContext().Scan.
func (b DeleteBuilder) ScanContext(ctx context.Context, dest ...interface{}) error {
	return b.QueryRowContext(ctx).Scan(dest...)
}

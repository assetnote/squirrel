// Package squirrel provides a fluent SQL generator.
//
// See https://github.com/assetnote/squirrel for examples.
package squirrel

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/lann/builder"
)

// Sqlizer is the interface that wraps the ToSql method.
//
// ToSql returns a SQL representation of the Sqlizer, along with a slice of args
// as passed to e.g. database/sql.Exec. It can also return an error.
type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}

// rawSqlizer is expected to do what Sqlizer does, but without finalizing placeholders.
// This is useful for nested queries.
type rawSqlizer interface {
	toSqlRaw() (string, []interface{}, error)
}

// Execer is the interface that wraps the Exec method.
//
// Exec executes the given query as implemented by database/sql.Exec.
type Execer interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
}

// Queryer is the interface that wraps the Query method.
//
// Query executes the given query as implemented by database/sql.Query.
type Queryer interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

// QueryRower is the interface that wraps the QueryRow method.
//
// QueryRow executes the given query as implemented by database/sql.QueryRow.
type QueryRower interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

// BaseRunner groups the Execer and Queryer interfaces.
type BaseRunner interface {
	Execer
	Queryer
}

// Runner groups the Execer, Queryer, and QueryRower interfaces.
type Runner interface {
	Execer
	Queryer
	QueryRower
}

// WrapStdSql wraps a type implementing the standard SQL interface with methods that
// squirrel expects.
func WrapStdSql(stdSql StdSql) Runner {
	return &stdsqlRunner{stdSql}
}

// StdPgxSql encompasses the standard methods of the pgx.Tx type, and other types that
// wrap these methods.
type StdSql interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type stdsqlRunner struct {
	StdSql
}

func (r *stdsqlRunner) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return r.StdSql.QueryRow(ctx, query, args...)
}

func setRunWith(b interface{}, runner BaseRunner) interface{} {
	switch r := runner.(type) {
	case StdSql:
		runner = WrapStdSql(r)
	}
	return builder.Set(b, "RunWith", runner)
}

// RunnerNotSet is returned by methods that need a Runner if it isn't set.
var RunnerNotSet = fmt.Errorf("cannot run; no Runner set (RunWith)")

// RunnerNotQueryRunner is returned by QueryRow if the RunWith value doesn't implement QueryRower.
var RunnerNotQueryRunner = fmt.Errorf("cannot QueryRow; Runner is not a QueryRower")

// ExecWith Execs the SQL returned by s with db.
func ExecWith(db Execer, s Sqlizer) (commandTag pgconn.CommandTag, err error) {
	query, args, err := s.ToSql()
	if err != nil {
		return
	}
	return db.Exec(context.Background(), query, args...)
}

// QueryWith Querys the SQL returned by s with db.
func QueryWith(db Queryer, s Sqlizer) (rows pgx.Rows, err error) {
	query, args, err := s.ToSql()
	if err != nil {
		return
	}
	return db.Query(context.Background(), query, args...)
}

// QueryRowWith QueryRows the SQL returned by s with db.
func QueryRowWith(db QueryRower, s Sqlizer) pgx.Row {
	query, args, err := s.ToSql()
	return &Row{Row: db.QueryRow(context.Background(), query, args...), err: err}
}

// DebugSqlizer calls ToSql on s and shows the approximate SQL to be executed
//
// If ToSql returns an error, the result of this method will look like:
// "[ToSql error: %s]" or "[DebugSqlizer error: %s]"
//
// IMPORTANT: As its name suggests, this function should only be used for
// debugging. While the string result *might* be valid SQL, this function does
// not try very hard to ensure it. Additionally, executing the output of this
// function with any untrusted user input is certainly insecure.
func DebugSqlizer(s Sqlizer) string {
	sql, args, err := s.ToSql()
	if err != nil {
		return fmt.Sprintf("[ToSql error: %s]", err)
	}

	var placeholder string
	downCast, ok := s.(placeholderDebugger)
	if !ok {
		placeholder = "?"
	} else {
		placeholder = downCast.debugPlaceholder()
	}
	// TODO: dedupe this with placeholder.go
	buf := &bytes.Buffer{}
	i := 0
	for {
		p := strings.Index(sql, placeholder)
		if p == -1 {
			break
		}
		if len(sql[p:]) > 1 && sql[p:p+2] == "??" { // escape ?? => ?
			buf.WriteString(sql[:p])
			buf.WriteString("?")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]
		} else {
			if i+1 > len(args) {
				return fmt.Sprintf(
					"[DebugSqlizer error: too many placeholders in %#v for %d args]",
					sql, len(args))
			}
			buf.WriteString(sql[:p])
			fmt.Fprintf(buf, "'%v'", args[i])
			// advance our sql string "cursor" beyond the arg we placed
			sql = sql[p+1:]
			i++
		}
	}
	if i < len(args) {
		return fmt.Sprintf(
			"[DebugSqlizer error: not enough placeholders in %#v for %d args]",
			sql, len(args))
	}
	// "append" any remaning sql that won't need interpolating
	buf.WriteString(sql)
	return buf.String()
}

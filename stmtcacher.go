package squirrel

// Deprecated: do no use. use pgx default stmt cache
//
// StmtCache was removed in the migration to pgx as the prepare semantics
// are different between database/sql and pgx
type StmtCache struct {
}


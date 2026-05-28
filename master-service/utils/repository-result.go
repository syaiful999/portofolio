package utils

import "database/sql"

type RepositoryResult struct {
	Result interface{}
	Error  error
	Sql    *sql.Rows
}

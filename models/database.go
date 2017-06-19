package models

import (
	"database/sql"
)

var DB *sql.DB

type Entity interface {
	TableName() string
	InsertString() string
}

type NotFoundError struct{}

func (e NotFoundError) Error() string {
	return "the record was not found"
}

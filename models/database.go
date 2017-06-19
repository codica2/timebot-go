package models

import (
	"database/sql"
)

var DB *sql.DB

type Entity interface {
	TableName() string
	InsertString() string
}

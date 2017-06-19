package models

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

type Entity interface {
	TableName() string
	InsertString() string
}

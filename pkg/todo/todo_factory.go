package todo

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var ErrUnknownStorageType = errors.New("unknown storage type")

type Config struct {
	StorageType string `json:"storage_type"`
	SQLDBPath    string `json:"sqldb_path"`
	HTTPPort      string `json:"http_port"`
}

func NewTodoServiceFromConfig(cfg Config) (TodoService, error) {
	switch cfg.StorageType {
	case "sqlite3":
		db, err := sql.Open("sqlite3", cfg.SQLDBPath)
		if err != nil {
			return nil, err
		}
		return NewTodoSQLite(db), nil
	case "mariadb":
		db, err := sql.Open("mysql", cfg.SQLDBPath)
		if err != nil {
			return nil, err
		}
		return NewTodoMariaDB(db), nil
	default:
		return nil, ErrUnknownStorageType
	}
}
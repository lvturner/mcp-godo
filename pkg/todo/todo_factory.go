package todo

import (
	"database/sql"
	"errors"
)

var ErrUnknownStorageType = errors.New("unknown storage type")

type Config struct {
	StorageType string `json:"storage_type"`
	SQLDBPath    string `json:"sqldb_path"`
	HTTPPort      string `json:"http_port"`
}

func NewTodoServiceFromConfig(cfg Config) (TodoService, error) {
	switch cfg.StorageType {

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

func NewProjectServiceFromConfig(cfg Config) (ProjectService, error) {
	switch cfg.StorageType {

	case "mariadb":
		db, err := sql.Open("mysql", cfg.SQLDBPath)
		if err != nil {
			return nil, err
		}
		return NewProjectMariaDB(db), nil
	default:
		return nil, ErrUnknownStorageType
	}
}

func NewCategoryServiceFromConfig(cfg Config) (CategoryService, error) {
	switch cfg.StorageType {

	case "mariadb":
		db, err := sql.Open("mysql", cfg.SQLDBPath)
		if err != nil {
			return nil, err
		}
		repo := NewCategoryMariaDB(db)
		return NewCategoryService(repo), nil
	default:
		return nil, ErrUnknownStorageType
	}
}
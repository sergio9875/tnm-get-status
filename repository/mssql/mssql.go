package mssql

import (
	"database/sql"
	"fmt"
	"malawi-getstatus/models"
	repo "malawi-getstatus/repository"

	_ "github.com/denisenkom/go-mssqldb"
)

// repository represent the repository model
type repository struct {
	db *sql.DB
}

// NewRepository will create a variable that represent the Repository struct
func NewRepository(dbConfig *models.MssqlConfig) (repo.Repository, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Port, dbConfig.Database)
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(3)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &repository{db}, nil
}

// Close attaches the provider and close the connection
func (r *repository) Close() error {
	return r.db.Close()
}

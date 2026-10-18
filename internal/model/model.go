package model

import (
	"context"
	"database/sql"
	"errors"
)

const (
	QueryStage = "query"
	ScanStage  = "scan"
	RowsStage  = "rows"
)

// DataError identifica em qual etapa do acesso ao banco ocorreu uma falha.
type DataError struct {
	Stage string
	Err   error
}

func (e *DataError) Error() string {
	return e.Err.Error()
}

func (e *DataError) Unwrap() error {
	return e.Err
}

// Model contém o acesso aos dados da aplicação.
type Model struct {
	db *sql.DB
}

func New(db *sql.DB) *Model {
	return &Model{db: db}
}

func validateDB(db *sql.DB) error {
	if db == nil {
		return errors.New("database is not initialized")
	}
	return nil
}

func (m *Model) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if err := validateDB(m.db); err != nil {
		return nil, &DataError{Stage: QueryStage, Err: err}
	}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, &DataError{Stage: QueryStage, Err: err}
	}
	return rows, nil
}

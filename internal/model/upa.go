package model

import (
	"context"
	"database/sql"
	"fmt"
)

// UPA representa os dados de uma unidade de pronto atendimento.
type UPA struct {
	ID           string
	Lat          float64
	Long         float64
	Title        string
	Phone        string
	Street       string
	Number       string
	Neighborhood string
	City         string
	State        string
	ZipCode      string
}

func (m *Model) ListUPAs(ctx context.Context) ([]UPA, error) {
	rows, err := m.query(ctx, `
		SELECT id, lat, long, title, phone,
		       street, number, neighborhood, city, state, zip_code
		FROM upas
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var upas []UPA
	for rows.Next() {
		var upa UPA
		var phone, number, neighborhood, zipCode sql.NullString

		if err := rows.Scan(
			&upa.ID, &upa.Lat, &upa.Long, &upa.Title, &phone,
			&upa.Street, &number, &neighborhood, &upa.City, &upa.State, &zipCode,
		); err != nil {
			return nil, &DataError{Stage: ScanStage, Err: err}
		}

		upa.Phone = phone.String
		upa.Number = number.String
		upa.Neighborhood = neighborhood.String
		upa.ZipCode = zipCode.String
		upas = append(upas, upa)
	}

	if err := rows.Err(); err != nil {
		return nil, &DataError{Stage: RowsStage, Err: fmt.Errorf("read UPA rows: %w", err)}
	}

	return upas, nil
}

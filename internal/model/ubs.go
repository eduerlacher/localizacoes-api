package model

import (
	"context"
	"database/sql"
	"fmt"
)

// UBS representa os dados de uma unidade básica de saúde.
type UBS struct {
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
	OpeningHours map[string]*Schedule
}

// Schedule guarda os horários opcionais de abertura e fechamento.
type Schedule struct {
	Open  *string
	Close *string
}

func (m *Model) ListUBSs(ctx context.Context) ([]UBS, error) {
	rows, err := m.query(ctx, `
		SELECT id, lat, long, title, phone,
		       street, number, neighborhood, city, state, zip_code,
		       monday_open, monday_close,
		       tuesday_open, tuesday_close,
		       wednesday_open, wednesday_close,
		       thursday_open, thursday_close,
		       friday_open, friday_close,
		       saturday_open, saturday_close,
		       sunday_open, sunday_close
		FROM ubss
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ubss := []UBS{}
	for rows.Next() {
		var ubs UBS
		var phone, number, neighborhood, zipCode sql.NullString
		var mondayOpen, mondayClose sql.NullString
		var tuesdayOpen, tuesdayClose sql.NullString
		var wednesdayOpen, wednesdayClose sql.NullString
		var thursdayOpen, thursdayClose sql.NullString
		var fridayOpen, fridayClose sql.NullString
		var saturdayOpen, saturdayClose sql.NullString
		var sundayOpen, sundayClose sql.NullString

		if err := rows.Scan(
			&ubs.ID, &ubs.Lat, &ubs.Long, &ubs.Title, &phone,
			&ubs.Street, &number, &neighborhood, &ubs.City, &ubs.State, &zipCode,
			&mondayOpen, &mondayClose,
			&tuesdayOpen, &tuesdayClose,
			&wednesdayOpen, &wednesdayClose,
			&thursdayOpen, &thursdayClose,
			&fridayOpen, &fridayClose,
			&saturdayOpen, &saturdayClose,
			&sundayOpen, &sundayClose,
		); err != nil {
			return nil, &DataError{Stage: ScanStage, Err: err}
		}

		ubs.Phone = phone.String
		ubs.Number = number.String
		ubs.Neighborhood = neighborhood.String
		ubs.ZipCode = zipCode.String
		ubs.OpeningHours = map[string]*Schedule{
			"monday":    {Open: nullableString(mondayOpen), Close: nullableString(mondayClose)},
			"tuesday":   {Open: nullableString(tuesdayOpen), Close: nullableString(tuesdayClose)},
			"wednesday": {Open: nullableString(wednesdayOpen), Close: nullableString(wednesdayClose)},
			"thursday":  {Open: nullableString(thursdayOpen), Close: nullableString(thursdayClose)},
			"friday":    {Open: nullableString(fridayOpen), Close: nullableString(fridayClose)},
			"saturday":  weekendSchedule(saturdayOpen, saturdayClose),
			"sunday":    weekendSchedule(sundayOpen, sundayClose),
		}
		ubss = append(ubss, ubs)
	}

	if err := rows.Err(); err != nil {
		return nil, &DataError{Stage: RowsStage, Err: fmt.Errorf("read UBS rows: %w", err)}
	}

	return ubss, nil
}

func nullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func weekendSchedule(open, close sql.NullString) *Schedule {
	if !open.Valid && !close.Valid {
		return nil
	}
	return &Schedule{Open: nullableString(open), Close: nullableString(close)}
}

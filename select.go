package main

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

type AccountRecords struct {
	Account
	Record
	Recorded    bool
	temperature sql.NullFloat64
	tempType    sql.NullInt32
	createdAt   sql.NullTime
	reason      sql.NullString
}

func selectRecords(tx *gorm.DB) (result []AccountRecords) {
	var err error
	var length int
	rows, err := tx.Select(`
		accounts.id, classes.name, accounts.number, accounts.name,
		role, record, account, 
		temperature, type, records.created_at, records.reason
		`).Count(&length).Rows()
	if err != nil {
		panic(err)
	}

	result = make([]AccountRecords, length)
	for i := 0; rows.Next(); i++ {
		var r AccountRecords
		err = rows.Scan(
			&r.Account.ID, &r.Class.Name, &r.Number, &r.Name,
			&r.Role, &r.Authority.Record, &r.Authority.Account,
			&r.temperature, &r.tempType, &r.createdAt, &r.reason,
		)
		if err != nil {
			panic(err)
		}
		if r.temperature.Valid {
			r.Recorded = true
			r.Record.Temperature = r.temperature.Float64
			r.Record.Type = TempType(r.tempType.Int32)
			r.Record.CreatedAt = r.createdAt.Time
			r.Record.Reason = r.reason.String
		}
		result[i] = r
	}
	return
}

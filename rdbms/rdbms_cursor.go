package rdbms

import (
	"database/sql"
	"github.com/eaciit/database/base"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
)

type Cursor struct {
	base.CursorBase
	rows *sql.Rows
}

func createError(title string, message string) error {
	return errorlib.Error(packageName, modCursor, title, message)
}

func (c *Cursor) validate() error {
	if c.QueryString == "" {
		createError("validate", "Invalid Query or Pipe Object")
	}

	return nil
}

func (c *Cursor) FetchAll(result interface{}, closeCursor bool) error {
	if e := c.validate(); e != nil {
		return createError("FetchAll", e.Error())
	}

	session := c.Connection.(*Connection).Sql
	rowRaw, e := session.Query(c.QueryString)
	c.rows = rowRaw

	if e != nil {
		return createError("FetchAll", e.Error())
	}

	columns, e := c.rows.Columns()

	if e != nil {
		return createError("FetchAll", e.Error())
	}

	rowDataRaw := make([]sql.RawBytes, len(columns))
	rowMemory := make([]interface{}, len(rowDataRaw))
	for i := range rowDataRaw {
		rowMemory[i] = &rowDataRaw[i]
	}

	rowAll := make([]toolkit.M, 0)

	defer c.rows.Close()

	for c.rows.Next() {
		e := c.rows.Scan(rowMemory...)

		if e != nil {
			return createError("FetchAll", e.Error())
		}

		rowData := toolkit.M{}

		for i, each := range rowDataRaw {
			value := "NULL"

			if each != nil {
				value = string(each)
			}

			rowData.Set(columns[i], value)
		}

		rowAll = append(rowAll, rowData)
	}

	*(result.(*[]toolkit.M)) = rowAll

	return c.rows.Err()
}

// func (c *Cursor) ResetFetch() error {
// }

// func (c *Cursor) FetchClose(result interface{}) (bool, error) {
// }

// func (c *Cursor) Fetch(result interface{}) (bool, error) {
// }

// func (c *Cursor) Count() int {
// }

// func (c *Cursor) Close() {
// }

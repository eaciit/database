package rdbms

import (
	"database/sql"
	"github.com/eaciit/database/base"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
)

type Cursor struct {
	base.CursorBase
	rows       *sql.Rows
	columns    []string
	rowMemory  []interface{}
	rowDataRaw []sql.RawBytes
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

func (c *Cursor) prepareFetch() error {
	if e := c.validate(); e != nil {
		return createError("prepareFetch", e.Error())
	}

	session := c.Connection.(*Connection).Sql
	rowRaw, e := session.Query(c.QueryString)

	if e != nil {
		return createError("prepareFetch", e.Error())
	} else {
		c.rows = rowRaw
	}

	columns, e := c.rows.Columns()

	if e != nil {
		return createError("FetchAll", e.Error())
	} else {
		c.columns = columns
	}

	rowDataRaw := make([]sql.RawBytes, len(columns))
	rowMemory := make([]interface{}, len(rowDataRaw))
	for i := range rowDataRaw {
		rowMemory[i] = &rowDataRaw[i]
	}
	c.rowDataRaw = rowDataRaw
	c.rowMemory = rowMemory

	return nil
}

func (c *Cursor) FetchAll(result interface{}, closeCursor bool) error {
	if e := c.prepareFetch(); e != nil {
		return e
	}

	rowAll := make([]toolkit.M, 0)
	defer c.FetchClose()

	for {
		rowData := toolkit.M{}

		if isNext, e := c.Fetch(&rowData); !isNext {
			if e != nil {
				return e
			}
			break
		}

		rowAll = append(rowAll, rowData)
	}

	*(result.(*[]toolkit.M)) = rowAll

	return c.rows.Err()
}

func (c *Cursor) Fetch(result interface{}) (bool, error) {
	if !c.rows.Next() {
		return false, nil
	}

	e := c.rows.Scan(c.rowMemory...)

	if e != nil {
		return false, createError("Fetch", e.Error())
	}

	rowData := toolkit.M{}

	for i, each := range c.rowDataRaw {
		value := "NULL"

		if each != nil {
			value = string(each)
		}

		rowData.Set(c.columns[i], value)
	}

	*(result.(*toolkit.M)) = rowData

	return true, nil
}

// func (c *Cursor) ResetFetch() error {
// }

// func (c *Cursor) Count() int {
// }

func (c *Cursor) Close() {
	c.rows.Close()
}

func (c *Cursor) FetchClose(result interface{}) (bool, error) {
	c.Close()
	return true, nil
}

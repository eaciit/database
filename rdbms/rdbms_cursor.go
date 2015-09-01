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

func (c *Cursor) doFetch(nCount int, result interface{}, callbackEach func(toolkit.M), closeCursor bool) (int, error) {
	rowAll := make([]toolkit.M, 0)

	if closeCursor {
		defer c.Close()
	}

	var i, j int

	for {
		if nCount != -1 && i >= nCount {
			break
		}

		i++

		rowData := toolkit.M{}

		if isNext, e := c.Fetch(&rowData); !isNext {
			if e != nil {
				return j, e
			}
			break
		}

		if callbackEach != nil {
			callbackEach(rowData)
		}

		rowAll = append(rowAll, rowData)

		j++
	}

	*(result.(*[]toolkit.M)) = rowAll

	return j, nil
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

func (c *Cursor) FetchN(nCount int, resultDataSet *base.DataSet, closeCursor bool) error {
	resultInside := []toolkit.M{}

	_, e := c.doFetch(nCount, &resultInside, func(each toolkit.M) {
		resultDataSet.Data = append(resultDataSet.Data, each)
	}, closeCursor)

	if e != nil {
		return e
	}

	return nil
}

func (c *Cursor) FetchAll(result interface{}, closeCursor bool) error {
	resultInside := []toolkit.M{}

	if _, e := c.doFetch(-1, &resultInside, nil, closeCursor); e != nil {
		return e
	}

	*(result.(*[]toolkit.M)) = resultInside

	return nil
}

func (c *Cursor) FetchClose(result interface{}) (bool, error) {
	resultInside := []toolkit.M{}

	if _, e := c.doFetch(-1, &resultInside, nil, true); e != nil {
		return false, e
	}

	*(result.(*[]toolkit.M)) = resultInside

	return true, nil
}

func (c *Cursor) ResetFetch() error {
	return c.prepareFetch()
}

func (c *Cursor) Count() int {
	session := c.Connection.(*Connection).Sql
	rows, e := session.Query(c.QueryString)

	if e != nil {
		return 0
	}

	var counter int

	for rows.Next() {
		counter++
	}

	return counter
}

func (c *Cursor) Close() {
	c.rows.Close()
}

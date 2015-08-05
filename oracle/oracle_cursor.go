package oracle

import (
	"github.com/eaciit/database/base"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
)

type Cursor struct {
	base.CursorBase
	QueryString string
}

func createError(title string, message string) error {
	return errorlib.Error(packageName, modCursor, "validate", "Invalid Query or Pipe Object")
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

	rowRaw, e := c.Connection.(*Connection).Sql.Query(c.QueryString)

	if e != nil {
		return createError("FetchAll", e.Error())
	}

	columns, e := rowRaw.Columns()

	if e != nil {
		return createError("FetchAll", e.Error())
	}

	allRows := make([]toolkit.M, 0)

	for rowRaw.Next() {
		var rowElement []interface{}
		rowElement = append(rowElement, make([]interface{}, (len(columns)-len(rowElement)))...)

		for i := range rowElement {
			rowElement[i] = &rowElement[i]
		}

		e := rowRaw.Scan(rowElement...)

		if e != nil {
			return createError("FetchAll", e.Error())
		}

		rowData := toolkit.M{}

		for i, each := range rowElement {
			rowData.Set(columns[i], each)
		}

		allRows = append(allRows, rowData)
	}

	*(result.(*[]toolkit.M)) = allRows

	return nil
}

package oracle

import (
	"github.com/eaciit/database/base"
	"github.com/eaciit/toolkit"
	"strings"
)

type Query struct {
	base.QueryBase
	currentParseMode string
}

func (q *Query) Parse(qe *base.QE, ins toolkit.M) interface{} {
	return qe.Value
}

func (q *Query) Compile(ins toolkit.M) (base.ICursor, interface{}, error) {
	s := q.Settings()
	qs := ""

	if s.Has("select") {
		qs += "select " + strings.Join(s.Get("select", []string{}).([]string), ",") + " "
	}

	if s.Has("from") {
		qs += "from " + s.Get("from", "").(string)
	}

	qs += " where customerid = 'ALFKI'"

	/* where clause in here */

	cursor := q.Connection.Table(qs, nil)

	return cursor, 0, nil
}

package oracle

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/toolkit"
	"strings"
)

type Query struct {
	base.QueryBase
	currentParseMode string
}

func (q *Query) parseWhere(op string, clauses []*base.QE) string {
	result := []string{}
	sep := ""

	if op == base.OpAnd {
		sep = "AND"
	} else if op == base.OpOr {
		sep = "OR"
	}

	for _, clause := range clauses {
		var clauseString string

		if clause.FieldOp == base.OpAnd || clause.FieldOp == base.OpOr {
			subWhere := q.parseWhere(clause.FieldOp, clause.Value.([]*base.QE))
			clauseString = fmt.Sprintf("(%s)", subWhere)
		} else if clause.FieldOp == base.OpEq {
			clauseString = fmt.Sprintf("%s = %s", clause.FieldId, clause.Value.(string))
		} else if clause.FieldOp == base.OpNe {
			clauseString = fmt.Sprintf("%s <> %s", clause.FieldId, clause.Value.(string))
		} else if clause.FieldOp == base.OpGt {
			clauseString = fmt.Sprintf("%s > %s", clause.FieldId, clause.Value.(string))
		} else if clause.FieldOp == base.OpGte {
			clauseString = fmt.Sprintf("%s >= %s", clause.FieldId, clause.Value.(string))
		} else if clause.FieldOp == base.OpLt {
			clauseString = fmt.Sprintf("%s < %s", clause.FieldId, clause.Value.(string))
		} else if clause.FieldOp == base.OpLte {
			clauseString = fmt.Sprintf("%s <= %s", clause.FieldId, clause.Value.(string))
		} else if clause.FieldOp == base.OpIn {
			value := strings.Join(clause.Value.([]string), ", ")
			clauseString = fmt.Sprintf("%s in (%s)", clause.FieldId, value)
		}

		result = append(result, clauseString)
	}

	return strings.Join(result, fmt.Sprintf(" %s ", sep))
}

func (q *Query) Parse(qe *base.QE, ins toolkit.M) interface{} {
	if qe.FieldOp == base.OpSelect {
		return qe.Value
	} else if qe.FieldOp == base.OpFromTable {
		return qe.Value
	} else if qe.FieldOp == base.OpAnd || qe.FieldOp == base.OpOr {
		parsedWhere := q.parseWhere(qe.FieldOp, qe.Value.([]*base.QE))
		for k, v := range ins {
			parsedWhere = strings.Replace(parsedWhere, k, q.StringValue(v.(string)), -1)
		}
		return parsedWhere
	} else if qe.FieldOp == base.OpOrderBy {
		parsedOrder := strings.Join(qe.Value.([]string), ", ")
		return parsedOrder
	}

	return qe.Value
}

func (q *Query) Compile(ins toolkit.M) (base.ICursor, interface{}, error) {
	s := q.Settings()
	qs := ""

	if s.Has("select") {
		queryPart := strings.Join(s.Get("select", []string{}).([]string), ", ")
		qs = fmt.Sprintf("%sSELECT %s ", qs, queryPart)
	}

	if s.Has("from") {
		queryPart := s.Get("from", "").(string)
		qs = fmt.Sprintf("%sFROM %s ", qs, queryPart)
	}

	if s.Has("where") {
		queryPart := s.Get("where", "").(string)
		qs = fmt.Sprintf("%sWHERE %s ", qs, queryPart)
	}

	if s.Has("orderby") {
		queryPart := s.Get("orderby", "").(string)
		qs = fmt.Sprintf("%sORDER BY %s ", qs, queryPart)
	}

	fmt.Println(qs)

	cursor := q.Connection.Table(qs, nil)

	return cursor, 0, nil
}

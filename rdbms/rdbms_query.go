package rdbms

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/toolkit"
	"strconv"
	"strings"
)

type Query struct {
	base.QueryBase
	currentParseMode string
}

func (q *Query) parseWhere(op string, clauses []*base.QE, ins toolkit.M) string {
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
			subWhere := q.parseWhere(clause.FieldOp, clause.Value.([]*base.QE), ins)
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
		} else if clause.FieldOp == base.OpContains {
			value := q.StringValue("")
			if len(clause.Value.([]string)) > 0 {
				value = fmt.Sprintf("%s%s%s", "%", ins.Get(clause.Value.([]string)[0], ""), "%")
			}
			clauseString = fmt.Sprintf("%s like %s", clause.FieldId, q.StringValue(value))
		} else if clause.FieldOp == base.OpStartWith {
			value := fmt.Sprintf("%s%s", ins.Get(clause.Value.(string), ""), "%")
			clauseString = fmt.Sprintf("%s like %s", clause.FieldId, q.StringValue(value))
		} else if clause.FieldOp == base.OpEndWith {
			value := fmt.Sprintf("%s%s", "%", ins.Get(clause.Value.(string), ""))
			clauseString = fmt.Sprintf("%s like %s", clause.FieldId, q.StringValue(value))
		} else if clause.FieldOp == base.OpBetween {
			value := []interface{}{
				clause.FieldId,
				clause.Value.(base.QRange).From,
				clause.FieldId,
				clause.Value.(base.QRange).To,
			}
			clauseString = fmt.Sprint("(", value[0], " >= ", value[1], " AND ", value[2], " <= ", value[3], ")")
		}

		result = append(result, clauseString)
	}

	parsedWhere := strings.Join(result, fmt.Sprintf(" %s ", sep))

	for k, v := range ins {
		var value string

		switch v.(type) {
		case int:
			value = strconv.Itoa(v.(int))
		default:
			value = q.StringValue(v.(string))
		}

		parsedWhere = strings.Replace(parsedWhere, k, value, -1)
	}

	return parsedWhere
}

func (q *Query) Parse(qe *base.QE, ins toolkit.M) interface{} {
	if qe.FieldOp == base.OpSelect {
		return qe.Value
	} else if qe.FieldOp == base.OpFromTable {
		return qe.Value
	} else if qe.FieldOp == base.OpAnd || qe.FieldOp == base.OpOr {
		parsedWhere := q.parseWhere(qe.FieldOp, qe.Value.([]*base.QE), ins)
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

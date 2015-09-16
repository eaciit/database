package rdbms

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/toolkit"
	// "os"
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
		value := q.getAsString(v)
		parsedWhere = strings.Replace(parsedWhere, k, value, -1)
	}

	return parsedWhere
}

func (q *Query) Parse(qe *base.QE, ins toolkit.M) interface{} {
	if qe.FieldOp == base.OpSelect {
		parsedSelect := strings.Join(qe.Value.([]string), ", ")
		return parsedSelect
	} else if qe.FieldOp == base.OpOrderBy {
		parsedOrder := strings.Join(qe.Value.([]string), ", ")
		return parsedOrder
	} else if qe.FieldOp == base.OpFromTable {
		parsedTable := strings.Join(qe.Value.([]string), ", ")
		return parsedTable
	} else if qe.FieldOp == base.OpWhereString || qe.FieldOp == base.OpAnd || qe.FieldOp == base.OpOr {
		if qe.FieldOp == base.OpWhereString {
			return qe.Value.(string)
		} else {
			parsedWhere := q.parseWhere(qe.FieldOp, qe.Value.([]*base.QE), ins)
			return parsedWhere
		}
	} else if qe.FieldOp == base.OpGroupBy {
		parsedGroup := strings.Join(qe.Value.([]string), ", ")
		return parsedGroup
	}

	return qe.Value
}

func (q *Query) Compile(ins toolkit.M) (base.ICursor, interface{}, error) {
	settings := q.Settings()
	commandType := q.CommandType(settings)
	queryString := ""

	compileNow := func(err ...error) (base.ICursor, interface{}, error) {
		cursor := q.Connection.Table(queryString, nil)
		errResetFetch := cursor.ResetFetch()

		if errResetFetch != nil {
			return nil, 0, errResetFetch
		}

		if len(err) > 0 {
			return nil, 0, err[0]
		}

		return cursor, 0, nil
	}

	if commandType == base.DB_SELECT {
		if settings.Has("select") {
			queryPart := settings.Get("select", "").(string)
			queryString = fmt.Sprintf("%sSELECT %s ", queryString, queryPart)
		}

		if settings.Has("from") {
			queryPart := settings.Get("from", "").(string)
			queryString = fmt.Sprintf("%sFROM %s ", queryString, queryPart)
		}

		if settings.Has("where") || settings.Has("whereString") {
			var whereKey = "where"
			if settings.Has("whereString") {
				whereKey = "whereString"
			}

			queryPart := settings.Get(whereKey, "").(string)
			queryString = fmt.Sprintf("%sWHERE %s ", queryString, queryPart)
		}

		if settings.Has("groupby") {
			queryPart := settings.Get("groupby", "").(string)
			queryString = fmt.Sprintf("%sGROUP BY %s ", queryString, queryPart)
		}

		if settings.Has("orderby") && q.Connection.(*Connection).Driver != "mssql" {
			queryPart := settings.Get("orderby", "").(string)
			queryString = fmt.Sprintf("%sORDER BY %s ", queryString, queryPart)
		}

		if q.Connection.(*Connection).Driver == "oci8" {
			queryString = q.compileLimitSkipForOracle(queryString)
		} else if q.Connection.(*Connection).Driver == "mssql" {
			queryString = q.compileLimitSkipForSqlServer(queryString)
		} else {
			if settings.Has("limit") {
				queryPart := settings.Get("limit", 10).(int)
				queryString = fmt.Sprintf("%sLIMIT %d ", queryString, queryPart)
			}

			if settings.Has("skip") {
				queryPart := settings.Get("skip", 0).(int)
				queryString = fmt.Sprintf("%sOFFSET %d ", queryString, queryPart)
			}
		}
	} else if commandType == base.DB_INSERT {
		if settings.Has("from") {
			queryPart := settings.Get("from", "").(string)
			queryString = fmt.Sprintf("INSERT INTO %s ", queryPart)
		} else {
			return compileNow(createError("Compile", "keyword FROM not found"))
		}

		queryString = q.compileInsertFrom(queryString, ins)
	} else if commandType == base.DB_UPDATE {
		if settings.Has("from") {
			queryPart := settings.Get("from", "").(string)
			queryString = fmt.Sprintf("UPDATE %s ", queryPart)
		} else {
			return compileNow(createError("Compile", "keyword FROM not found"))
		}

		queryString = q.compileUpdateFrom(queryString, ins)

		if settings.Has("where") {
			queryPart := settings.Get("where", "").(string)
			queryString = fmt.Sprintf("%sWHERE %s ", queryString, queryPart)
		}

		queryString = q.compileWhereBinding(queryString, ins)
	} else if commandType == base.DB_DELETE {
		if settings.Has("from") {
			queryPart := settings.Get("from", "").(string)
			queryString = fmt.Sprintf("%sDELETE FROM %s ", queryString, queryPart)
		}

		if settings.Has("where") {
			queryPart := settings.Get("where", "").(string)
			queryString = fmt.Sprintf("%sWHERE %s ", queryString, queryPart)
		}

		queryString = q.compileWhereBinding(queryString, ins)
	}
	fmt.Println(queryString)
	return compileNow()
}

func (q *Query) compileLimitSkipForOracle(queryString string) string {
	settings := q.Settings()

	if settings.Has("limit") && settings.Has("skip") {
		querySelect := settings.Get("select", "").(string)
		queryLimit := settings.Get("limit", 10).(int)
		querySkip := settings.Get("skip", 10).(int)
		queryString = fmt.Sprintf("SELECT %s FROM (SELECT table_inner.*, ROWNUM as table_counter from (%s) table_inner) WHERE table_counter > %d and (table_counter - %d) <= %d", querySelect, queryString, querySkip, querySkip, queryLimit)
	} else if settings.Has("limit") {
		querySelect := settings.Get("select", "").(string)
		queryPart := settings.Get("limit", 10).(int)
		queryString = fmt.Sprintf("SELECT %s FROM (SELECT table_inner.*, ROWNUM as table_counter from (%s) table_inner) WHERE table_counter <= %d", querySelect, queryString, queryPart)
	} else if settings.Has("skip") {
		querySelect := settings.Get("select", "").(string)
		queryPart := settings.Get("skip", 10).(int)
		queryString = fmt.Sprintf("SELECT %s FROM (SELECT table_inner.*, ROWNUM as table_counter from (%s) table_inner) WHERE table_counter > %d", querySelect, queryString, queryPart)
	}

	// e := createError("Compile", "Limit & Offset currently is not support on oracle driver")
	// fmt.Println(e.Error())
	// os.Exit(0)

	return queryString
}

func (q *Query) compileLimitSkipForSqlServer(queryString string) string {
	settings := q.Settings()
	rowNumber := "ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) as table_counter"

	if settings.Has("orderby") {
		queryPart := settings.Get("orderby", "").(string)
		rowNumber = strings.Replace(rowNumber, "(SELECT NULL)", queryPart, -1)
	}

	queryString = strings.Replace(queryString, " FROM", fmt.Sprintf(", %s FROM", rowNumber), -1)

	if settings.Has("limit") && settings.Has("skip") {
		querySelect := settings.Get("select", "").(string)
		queryLimit := settings.Get("limit", 10).(int)
		querySkip := settings.Get("skip", 10).(int)
		queryString = fmt.Sprintf("SELECT TOP %d %s FROM (%s) table_inner WHERE table_counter > %d", queryLimit, querySelect, queryString, querySkip)
	} else if settings.Has("limit") {
		querySelect := settings.Get("select", "").(string)
		queryPart := settings.Get("limit", 10).(int)
		queryString = fmt.Sprintf("SELECT TOP %d %s FROM (%s) table_inner", queryPart, querySelect, queryString)
	} else if settings.Has("skip") {
		querySelect := settings.Get("select", "").(string)
		queryPart := settings.Get("skip", 10).(int)
		queryString = fmt.Sprintf("SELECT %s FROM (%s) table_inner WHERE table_counter > %d", querySelect, queryString, queryPart)
	}

	return queryString
}

func (q *Query) compileInsertFrom(queryString string, ins toolkit.M) string {
	keyString, valString := func() (string, string) {
		if !ins.Has("data") {
			return "", ""
		}

		keys := []string{}
		vals := []string{}

		for k, v := range ins.Get("data", toolkit.M{}).(toolkit.M) {
			keys = append(keys, k)
			vals = append(vals, q.getAsString(v))
		}

		return strings.Join(keys, ", "), strings.Join(vals, ", ")
	}()

	queryString = fmt.Sprintf("%s(%s) VALUES (%s) ", queryString, keyString, valString)
	return queryString
}

func (q *Query) compileUpdateFrom(queryString string, ins toolkit.M) string {
	updateString := func() string {
		if !ins.Has("data") {
			return ""
		}

		var updates []string

		for k, v := range ins.Get("data", toolkit.M{}).(toolkit.M) {
			updates = append(updates, fmt.Sprintf("%s = %s", k, q.getAsString(v)))
		}

		return strings.Join(updates, ", ")
	}()

	queryString = fmt.Sprintf("%sSET %s ", queryString, updateString)
	return queryString
}

func (q *Query) compileWhereBinding(queryString string, ins toolkit.M) string {
	for k, v := range ins {
		if k == "data" {
			continue
		}

		queryString = strings.Replace(queryString, k, q.getAsString(v), -1)
	}

	return queryString
}

func (q *Query) getAsString(v interface{}) string {
	value := ""

	switch v.(type) {
	case int:
		value = strconv.Itoa(v.(int))
	case string:
		value = q.StringValue(v.(string))
	default:
		value = fmt.Sprintf("%v", v)
	}

	return value
}

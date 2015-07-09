package query

/******************************************************************************************
How to use:
Eq("Field1","1").And().Eq("Field2","2").Club().Or()Eq("Field1","root").ToWhere()
******************************************************************************************/

import (
	"fmt"
	"github.com/eaciit/toolkit"
)

const (
	OpEq           = "$eq"
	OpNe           = "$ne"
	OpGt           = "$gt"
	OpGte          = "$gte"
	OpLte          = "$lte"
	OpBetween      = "$between"
	OpIn           = "$in"
	OpContains     = "$contains"
	OpStartWith    = "$startwith"
	OpEndWith      = "$endwith"
	OpOpenBracket  = "$("
	OpCloseBracket = "$)"
	OpAnd          = "$and"
	OpOr           = "$or"
)

type QE struct {
	FieldId string
	FieldOp string
	Value   interface{}
}

type QRange struct {
	From interface{}
	To   interface{}
}

type IQuery interface {
	Eq(string, interface{}) IQuery
	Ne(string, interface{}) IQuery
	Gt(string, interface{}) IQuery
	Gte(string, interface{}) IQuery
	Lt(string, interface{}) IQuery
	Lte(string, interface{}) IQuery
	Contains(string, interface{}) IQuery
	StartWith(string, string) IQuery
	EndWith(string, string) IQuery
	Between(string, interface{}, interface{}) IQuery
	In(string, ...interface{}) IQuery
	SetStringSign(string) IQuery
	O() IQuery
	C() IQuery
	And() IQuery
	Or() IQuery
	Parse(toolkit.M) interface{}
}

type Query struct {
	elements   []QE
	stringSign string
}

func (q *Query) SetStringSign(str string) IQuery {
	q.stringSign = str
	return q
}

func (q *Query) add(qe IQueryElement) IQuery {
	q.elements = append(q.elements, qe)
	return q
}

func (q *Query) Eq(field string, value interface{}) IQuery {
	q.add(&QE{field, OpEq, value})
	return q
}

func (q *Query) Ne(field string, value interface{}) IQuery {
	q.add(&QE{field, OpNe, value})
	return q
}

func (q *Query) Gt(field string, value interface{}) IQuery {
	q.add(&QE{field, OpGt, value})
	return q
}

func (q *Query) Gte(field string, value interface{}) IQuery {
	q.add(&QE{field, OpGte, value})
	return q
}

func (q *Query) Lt(field string, value interface{}) IQuery {
	q.add(&QE{field, OpLt, value})
	return q
}

func (q *Query) Lte(field string, value interface{}) IQuery {
	q.add(&QE{field, OpLte, value})
	return q
}

func (q *Query) Contains(field string, value string) IQuery {
	q.add(&QE{field, OpContains, value})
	return q
}

func (q *Query) StartWith(field string, value string) IQuery {
	q.add(&QE{field, OpStartWith, value})
	return q
}

func (q *Query) EndWith(field string, value string) IQuery {
	q.add(&QE{field, OpEndWith, value})
	return q
}

func (q *Query) Between(field string, from interface{}, to interface{}) IQuery {
	q.add(&QE{field, OpBetween, QRange{from, to}})
	return q
}

func (q *Query) In(fields ...interface{}) IQuery {
	slices := make([]interface{}, 0)
	for _, v := range slices {
		slices = append(v)
	}
	q.add(&QE{field, OpInRange, slices})
	return q
}

func (q *Query) And() IQuery {
	q.add(&QE{"", OpAnd, nil})
	return q
}

func (q *Query) Or() IQuery {
	q.add(&QE{"", OpOr, nil})
	return q
}

func (q *Query) O() IQuery {
	q.add(&QE{"", OpOpenBracket, nil})
	return q
}

func (q *Query) C() IQuery {
	q.add(&QE{"", OpCloseBracket, nil})
	return q
}

func (q *Query) ParseValue(v interface{}) string {
	var ret string
	switch v.Type() {
	case string:
		if q.stringSign == "'" {
			ret = fmt.Sprintf("'%s'", v.(string))
		} else if q.stringSign == "\"" {
			ret = fmt.Sprintf("\"%s\"", v.(string))
		} else {
			ret = fmt.Sprintf("%s%s%s", q.stringSign v.(string), q.stringSign)
		}
		break
		
	case nil:
		ret = ""

	default:
		ret = fmt.Sprintf("%v", v)
		//-- do nothing
	}
	return ret
}

func (q *Query) Parse(ins toolkit.M) interface{} {
	part := ""
	command := ""

	for _, v := range q.elements {
		if v.FieldOp == OpOpenBracket {
			part = "("
		} else if v.FieldOp == OpCloseBracket {
			part = part + ")"
			command = command + part
		} else if v.FieldOp == OpEq {
			part = part + fmt.Sprintf("%s = %s", v.FieldId, v.Value)
		}
	}
}

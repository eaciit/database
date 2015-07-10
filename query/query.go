package query

/******************************************************************************************
How to use:
Eq("Field1","1").And().Eq("Field2","2").Club().Or()Eq("Field1","root").ToWhere()
******************************************************************************************/

import (
	"fmt"
	"github.com/eaciit/toolkit"
	"time"
)

const (
	OpEq           = "$eq"
	OpNe           = "$ne"
	OpGt           = "$gt"
	OpGte          = "$gte"
	OpLt           = "$lt"
	OpLte          = "$lte"
	OpLt           = "$lt"
	OpBetween      = "$between"
	OpIn           = "$in"
	OpContains     = "$contains"
	OpStartWith    = "$startwith"
	OpEndWith      = "$endwith"
	OpOpenBracket  = "$("
	OpCloseBracket = "$)"
	OpAnd          = "$and"
	OpOr           = "$or"
	OpTable        = "$table"
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
	Contains(string, ...string) IQuery
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

	SetQ(IQuery)
	Q() IQuery
}

type Query struct {
	Elements   []*QE
	q          IQuery
	stringSign string
}

func NewQuery(q IQuery) IQuery {
	q.SetQ(q)
	return q
}

func (q *Query) SetQ(self IQuery) {
	q.q = self
}

func (q *Query) Q() IQuery {
	return q.q
}

func (q *Query) SetStringSign(str string) IQuery {
	q.stringSign = str
	return q.Q()
}

func (q *Query) add(qe *QE) IQuery {
	q.Elements = append(q.Elements, qe)
	return q
}

func (q *Query) Eq(field string, value interface{}) IQuery {
	q.add(&QE{field, OpEq, value})
	return q.Q()
}

func (q *Query) Ne(field string, value interface{}) IQuery {
	q.add(&QE{field, OpNe, value})
	return q.Q()
}

func (q *Query) Gt(field string, value interface{}) IQuery {
	q.add(&QE{field, OpGt, value})
	return q.Q()
}

func (q *Query) Gte(field string, value interface{}) IQuery {
	q.add(&QE{field, OpGte, value})
	return q.Q()
}

func (q *Query) Lt(field string, value interface{}) IQuery {
	q.add(&QE{field, OpLt, value})
	return q.Q()
}

func (q *Query) Lte(field string, value interface{}) IQuery {
	q.add(&QE{field, OpLte, value})
	return q.Q()
}

func (q *Query) Contains(field string, value ...string) IQuery {
	q.add(&QE{field, OpContains, value})
	return q.Q()
}

func (q *Query) StartWith(field string, value string) IQuery {
	q.add(&QE{field, OpStartWith, value})
	return q.Q()
}

func (q *Query) EndWith(field string, value string) IQuery {
	q.add(&QE{field, OpEndWith, value})
	return q.Q()
}

func (q *Query) Between(field string, from interface{}, to interface{}) IQuery {
	q.add(&QE{field, OpBetween, QRange{from, to}})
	return q.Q()
}

func (q *Query) In(field string, slices ...interface{}) IQuery {
	q.add(&QE{field, OpIn, slices})
	return q.Q()
}

func (q *Query) And() IQuery {
	q.add(&QE{"", OpAnd, nil})
	return q.Q()
}

func (q *Query) Or() IQuery {
	q.add(&QE{"", OpOr, nil})
	return q.Q()
}

func (q *Query) O() IQuery {
	q.add(&QE{"", OpOpenBracket, nil})
	return q.Q()
}

func (q *Query) C() IQuery {
	q.add(&QE{"", OpCloseBracket, nil})
	return q.Q()
}

func (q *Query) ParseValue(v interface{}) string {
	var ret string
	switch v.(type) {
	case string:
		if q.stringSign == "'" {
			ret = fmt.Sprintf("'%s'", v.(string))
		} else if q.stringSign == "\"" {
			ret = fmt.Sprintf("\"%s\"", v.(string))
		} else {
			ret = fmt.Sprintf("%s%s%s", q.stringSign, v.(string), q.stringSign)
		}
	case time.Time:
		ret = fmt.Sprintf("%s%v%s", q.stringSign, v.(time.Time).UTC(), q.stringSign)

	case *time.Time:
		ret = fmt.Sprintf("%s%v%s", q.stringSign, v.(*time.Time), q.stringSign)

	case int, int32, int64, uint, uint32, uint64:
		ret = fmt.Sprintf("%d", v.(int))

	case nil:
		ret = ""

	default:
		ret = fmt.Sprintf("%v", v)
		//-- do nothing
	}
	return ret
}

func (q *Query) Parse(ins toolkit.M) interface{} {
	return nil
	part := ""
	command := ""

	for _, v := range q.Elements {
		if v.FieldOp == OpOpenBracket {
			command = command + part
			part = "("
		} else if v.FieldOp == OpCloseBracket {
			part = part + ")"
			command = command + part
		} else if v.FieldOp == OpOr {
			part = part + " or "
		} else if v.FieldOp == OpAnd {
			part = part + " and "
		} else if v.FieldOp == OpEq {
			part = part + fmt.Sprintf("%s=%s", v.FieldId, q.ParseValue(v.Value))
		}
	}
	return command
}

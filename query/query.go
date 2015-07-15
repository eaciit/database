package query

/******************************************************************************************
How to use:
Eq("Field1","1").And().Eq("Field2","2").Club().Or()Eq("Field1","root").ToWhere()
******************************************************************************************/

import (
	"fmt"
	. "github.com/eaciit/toolkit"
	"time"
)

const (
	OpEq           = "$eq"
	OpNe           = "$ne"
	OpGt           = "$gt"
	OpGte          = "$gte"
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
	OpChain        = "$chain"
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
	SetStringSign(string) IQuery
	Command(*M, *M, ...QE) error
	StringValue(interface{}) string
	Parse(QE, *M) interface{}

	SetQ(IQuery)
	Q() IQuery
}

type Query struct {
	Elements   []QE
	q          IQuery
	stringSign string
}

func New(q IQuery) IQuery {
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
	return q
}

func Eq(field string, value interface{}) QE {
	return QE{field, OpEq, value}
}

func Ne(field string, value interface{}) QE {
	return QE{field, OpNe, value}

}

func Gt(field string, value interface{}) QE {
	return QE{field, OpGt, value}

}

func Gte(field string, value interface{}) QE {
	return QE{field, OpGte, value}

}

func Lt(field string, value interface{}) QE {
	return QE{field, OpLt, value}

}

func Lte(field string, value interface{}) QE {
	return QE{field, OpLte, value}

}

func Contains(field string, value ...string) QE {
	return QE{field, OpContains, value}

}

func StartWith(field string, value string) QE {
	return QE{field, OpStartWith, value}

}

func EndWith(field string, value string) QE {
	return QE{field, OpEndWith, value}

}

func Between(field string, from interface{}, to interface{}) QE {
	return QE{field, OpBetween, QRange{from, to}}
}

func In(field string, invalues ...interface{}) QE {
	return QE{field, OpIn, invalues}

}

func And(qes ...QE) QE {
	return QE{"", OpAnd, qes}

}

func Or(qes ...QE) QE {
	return QE{"", OpOr, qes}

}

/*func (q *Query) O() QE {
	return QE{"", OpOpenBracket, nil}

}

func (q *Query) C() QE {
	return QE{"", OpCloseBracket, nil}

}

func (q *Query) Chain(chainQuery IQuery) QE {
	return QE{"", OpChain, chainQuery}

}
*/

func (q *Query) Command(result *M, ins *M, qes ...QE) error {
	m := *result
	if !m.Has("Data") {
		m.Set("Data", "")
	}
	if len(qes) == 1 {
		m.Set("Data", q.Parse(qes[0], ins))
	} else if len(qes) > 1 {
		newqs := And(qes...)
		m.Set("Data", q.Parse(newqs, ins))
	}
	return nil
}

func (q *Query) StringValue(v interface{}) string {
	var ret string
	switch v.(type) {
	case string:
		if q.stringSign == "'" {
			ret = fmt.Sprintf("'%s'", v.(string))
		} else if q.stringSign == "\"" {
			ret = fmt.Sprintf("\"%s\"", v.(string))
		} else if q.stringSign == "" {
			ret = fmt.Sprintf("'%s'", v.(string))
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

func (q *Query) Parse(qe QE, ins *M) interface{} {
	var v QE
	result := ""

	if qe.FieldOp == OpEq {
		result = fmt.Sprintf("%s = %s", qe.FieldId, q.StringValue(qe.Value))
	} else if qe.FieldOp == OpNe {
		result = fmt.Sprintf("%s != %s", qe.FieldId, q.StringValue(qe.Value))
	} else if qe.FieldOp == OpOr {
		tmp := ""
		for _, v = range qe.Value.([]QE) {
			if tmp != "" {
				tmp = tmp + " or "
			}
			tmp = tmp + q.Parse(v, ins).(string)
		}
		result = "(" + tmp + ")"
	} else if qe.FieldOp == OpAnd {
		tmp := ""
		for _, v = range qe.Value.([]QE) {
			if tmp != "" {
				tmp = tmp + " and "
			}
			tmp = tmp + q.Parse(v, ins).(string)
		}
		result = "(" + tmp + ")"
	}

	return result
}

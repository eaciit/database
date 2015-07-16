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
	Command(*Result, *M, ...QE) error
	StringValue(interface{}) string
	Parse(QE, *M) interface{}

	SetStringSign(string) IQuery
	SetQ(IQuery) IQuery
	Q() IQuery
}

type QueryBase struct {
	stringSign string
	q          IQuery
}

func (q *QueryBase) SetQ(i IQuery) IQuery {
	q.q = i
	return q
}

func (q *QueryBase) Q() IQuery {
	return q.q
}

func (q *QueryBase) SetStringSign(str string) IQuery {
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

func (q *QueryBase) Command(result *Result, ins *M, qes ...QE) error {
	if q.q == nil {
		result.Status = Status_NOK
		result.Message = "Query object is not properly initiated. Please call SetQ"
		return fmt.Errorf("Query object is not properly initiated. Please call SetQ")
	}
	if len(qes) == 1 {
		result.Data = q.Q().Parse(qes[0], ins)
	} else if len(qes) > 1 {
		newqs := And(qes...)
		result.Data = q.Q().Parse(newqs, ins)
	}
	return nil
}

func (q *QueryBase) StringValue(v interface{}) string {
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

func (q *QueryBase) Parse(qe QE, ins *M) interface{} {
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

func New(q IQuery) IQuery {
	q.SetQ(q)
	return q
}

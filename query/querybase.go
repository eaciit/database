package query

import (
	"fmt"
	. "github.com/eaciit/toolkit"
	"time"
)

type IQuery interface {
	Build(*Result, *M)
	Compile(*M) (interface{}, error)
	StringValue(interface{}) string
	Parse(*QE, *M) interface{}

	Select(...string) IQuery
	Where(...*QE) IQuery
	OrderBy(*QE) IQuery
	GroupBy(*QE) IQuery
	From(string) IQuery
	Limit(int) IQuery
	Skip(int) IQuery
	//Command(string, *QE) IQuery

	SetStringSign(string) IQuery
	SetQ(IQuery) IQuery
	Q() IQuery
}

type QueryBase struct {
	stringSign string
	q          IQuery

	Elements map[string]*QE
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

func (q *QueryBase) addQE(key string, v *QE) {
	if q.Elements == nil {
		q.Elements = make(map[string]*QE)
	}
	q.Elements[key] = v
}

func (q *QueryBase) Select(fields ...string) IQuery {
	q.addQE("select", &QE{"", OpSelect, fields})
	return q
}

func (q *QueryBase) From(tablename string) IQuery {
	q.addQE("from", &QE{"", OpFromTable, tablename})
	return q
}

func (q *QueryBase) Where(qes ...*QE) IQuery {
	if len(qes) == 1 {
		q.addQE("where", qes[0])
		//result.Data = q.Q().Parse(qes[0], ins)
	} else if len(qes) > 1 {
		newqs := And(qes...)
		q.addQE("where", newqs)
	}
	return q
}

func (q *QueryBase) OrderBy(qe *QE) IQuery {
	q.addQE("orderby", qe)
	return q
}

func (q *QueryBase) GroupBy(qe *QE) IQuery {
	q.addQE("groupby", qe)
	return q
}

func (q *QueryBase) Skip(s int) IQuery {
	q.addQE("skip", &QE{"", OpSkip, s})
	return q
}

func (q *QueryBase) Limit(l int) IQuery {
	q.addQE("limit", &QE{"", OpLimit, l})
	return q
}

func (q *QueryBase) Build(result *Result, ins *M) {
	if q.q == nil {
		result.Status = Status_NOK
		result.Message = "Query object is not properly initiated. Please call SetQ"
	}

	m := M{}
	for k, v := range q.Elements {
		m[k] = q.Q().Parse(v, ins)
	}
	t, e := q.Q().Compile(&m)
	if e != nil {
		result.Status = Status_NOK
		result.Message = e.Error()
	}

	result.Status = Status_OK
	result.Data = t
}

func (q *QueryBase) Compile(ins *M) (interface{}, error) {
	ret := ""

	concat := func(s string, as ...string) string {
		for _, a := range as {
			s += " " + a
		}
		return s
	}

	if ins.Has("select") {
		ret = concat(ret, ins.Get("select", "").(string))
	}
	if ins.Has("from") {
		ret = concat(ret, ins.Get("from", "").(string))
	}
	if ins.Has("where") {
		ret = concat(ret, "where", ins.Get("where", "").(string))
	}
	return ret, nil
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

func (q *QueryBase) Parse(qe *QE, ins *M) interface{} {
	var v *QE
	result := ""

	if qe.FieldOp == OpSelect {
		selecteds := ""
		for _, f := range qe.Value.([]string) {
			if selecteds != "" {
				selecteds += ","
			}
			selecteds += f
		}
		result = fmt.Sprintf("select %s", selecteds)
	} else
	// handle from
	if qe.FieldOp == OpFromTable {
		result = fmt.Sprintf("from %s", qe.Value)
	} else
	// this is for order
	if qe.FieldOp == OpEq {
		result = fmt.Sprintf("%s = %s", qe.FieldId, q.StringValue(qe.Value))
	} else if qe.FieldOp == OpNe {
		result = fmt.Sprintf("%s != %s", qe.FieldId, q.StringValue(qe.Value))
	} else if qe.FieldOp == OpOr {
		tmp := ""
		for _, v = range qe.Value.([]*QE) {
			if tmp != "" {
				tmp = tmp + " or "
			}
			tmp = tmp + q.Parse(v, ins).(string)
		}
		result = "(" + tmp + ")"
	} else if qe.FieldOp == OpAnd {
		tmp := ""
		for _, v = range qe.Value.([]*QE) {
			if tmp != "" {
				tmp = tmp + " and "
			}
			tmp = tmp + q.Parse(v, ins).(string)
		}
		result = "(" + tmp + ")"
	}

	return result
}

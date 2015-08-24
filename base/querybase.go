package base

import (
	"fmt"
	"github.com/eaciit/errorlib"
	. "github.com/eaciit/toolkit"
	"reflect"
	"strings"
	"time"
)

type IQuery interface {
	Build(M) error
	Run(M) (ICursor, interface{}, error)
	Compile(M) (ICursor, interface{}, error)
	StringValue(interface{}) string
	Parse(*QE, M) interface{}

	Command(M) IQuery
	Cursor(M) ICursor
	Select(...string) IQuery
	SetFields(...string) IQuery
	Where(...*QE) IQuery
	OrderBy(...string) IQuery
	GroupBy(...string) IQuery
	Aggregate(...*QE) IQuery
	From(...string) IQuery
	Limit(int) IQuery
	Skip(int) IQuery
	//Command(string, *QE) IQuery

	//Transaction
	Insert() IQuery
	Save() IQuery
	Update() IQuery
	Delete() IQuery

	CommandType(M) DB_OP
	SetStringSign(string) IQuery
	SetQ(IQuery) IQuery
	SetConnection(IConnection) IQuery
	Q() IQuery
	Settings() M
	Reset()
}

type QueryBase struct {
	stringSign string
	q          IQuery
	Connection IConnection

	Elements map[string]*QE
	settings M
}

func (q *QueryBase) SetConnection(c IConnection) IQuery {
	q.Connection = c
	return q.Q()
}

func (q *QueryBase) SetQ(i IQuery) IQuery {
	q.q = i
	return q
}

func (q *QueryBase) Q() IQuery {
	return q.q
}

func (q *QueryBase) Settings() M {
	return q.settings
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

func (q *QueryBase) SetFields(fields ...string) IQuery {
	q.addQE("set", &QE{"", OpSetfield, fields})
	return q
}

func (q *QueryBase) Aggregate(aggregates ...*QE) IQuery {
	q.addQE("aggregate", &QE{"", OpAggregate, aggregates})
	return q
}

func (q *QueryBase) From(tablenames ...string) IQuery {
	q.addQE("from", &QE{"", OpFromTable, tablenames})
	return q
}

func (q *QueryBase) Where(qes ...*QE) IQuery {
	// return as array even only one query
	// because this will become a problem
	// for rdbms
	q.addQE("where", And(qes...))
	return q
}

func (q *QueryBase) OrderBy(fields ...string) IQuery {
	q.addQE("orderby", &QE{"", OpOrderBy, fields})
	return q
}

func (q *QueryBase) GroupBy(gs ...string) IQuery {
	q.addQE("groupby", &QE{"", OpGroupBy, gs})
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

func (q *QueryBase) Insert() IQuery {
	q.addQE("insert", &QE{})
	return q
}

func (q *QueryBase) Save() IQuery {
	q.addQE("save", &QE{})
	return q
}

func (q *QueryBase) Update() IQuery {
	q.addQE("update", &QE{})
	return q
}

func (q *QueryBase) Delete() IQuery {
	q.addQE("delete", &QE{})
	return q
}

// To add command pass a M object with following signature
//
// M["command"] => name of command
// M[whatever] => whatever parameter require to support the command, ie:
// M{"command":"aggregate", "pipe":pipes}
func (q *QueryBase) Command(ins M) IQuery {
	cmdname := ins.Get("command", "").(string)
	if cmdname != "" {
		cmdname = "cmd." + cmdname
		q.addQE(cmdname, &QE{cmdname, OpCommand, ins})
	}
	return q
}

func (q *QueryBase) CommandType(ins M) DB_OP {
	dbopType := DB_SELECT
	if ins.Has("select") {
		dbopType = DB_SELECT
	}
	if ins.Has("update") {
		dbopType = DB_UPDATE
	}
	if ins.Has("save") {
		dbopType = DB_SAVE
	}
	if ins.Has("insert") {
		dbopType = DB_INSERT
	}
	if ins.Has("delete") {
		dbopType = DB_DELETE
	}
	if ins.Has("command") {
		dbopType = DB_COMMAND
	}
	return dbopType
}

func (q *QueryBase) Reset() {
	q.settings = nil
}

func (q *QueryBase) Build(ins M) error {
	if q.q == nil {
		return errorlib.Error(packageName, modQuery, "Build", "Query object is not properly initiated. Please call SetQ")
	}

	if q.settings == nil {
		m := M{}
		for k, v := range q.Elements {
			m[k] = q.Q().Parse(v, ins)
		}
		//_ = "breakpoint"
		q.settings = m
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
			ret = fmt.Sprintf("%s", v.(string))
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

func (q *QueryBase) ParseValue(o interface{}, m M) interface{} {
	//_ = "breakpoint"
	ref := reflect.ValueOf(o)
	if !ref.IsValid() {
		return nil
	}
	if ref.Kind() == reflect.String && strings.HasPrefix(ref.String(), "@") {
		parmName := ref.String()
		if m.Has(parmName) {
			ret := m[parmName]
			rv := reflect.ValueOf(ret)
			if rv.Kind() == reflect.String {
				ret = q.q.StringValue(ret)
			}
			return ret
		} else {
			return nil
		}
	}
	return o
}

func (q *QueryBase) Parse(qe *QE, ins M) interface{} {
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

// Run is used to finalize query, it will return 3 objects: cursor, interface{} and error
// an M object is optional to be passed as entry
// m["data"] => any data that will be saved, normally for non select
// other than data will be used as parameters replacement
func (q *QueryBase) Run(ins M) (ICursor, interface{}, error) {
	var e error
	e = q.q.Build(ins)
	if e != nil {
		return nil, 0, errorlib.Error(packageName, modQuery, "Run", e.Error())
	}
	return q.q.Compile(ins)
}

func (q *QueryBase) Compile(ins M) (ICursor, interface{}, error) {
	return nil, nil, errorlib.Error(packageName, modQuery, "Compile", errorlib.NotYetImplemented)
}

func (q *QueryBase) Cursor(params M) ICursor {
	//_ = "breakpoint"
	c, _, _ := q.q.Run(params)
	return c
}

package mongodb

import (
	_ "fmt"
	. "github.com/eaciit/database/base"
	"github.com/eaciit/errorlib"
	. "github.com/eaciit/toolkit"
)

type Query struct {
	QueryBase
	currentParseMode string
}

func (q *Query) Compile(ins M) (ICommand, error) {
	commandType := q.CommandType(ins)
	command := new(Command)
	parm := M{}

	//_ = "breakpoint"
	tableName := ""
	if ins.Has("from") {
		tableName = ins.Get("from", "").(string)
	} else {
		return nil, errorlib.Error(packageName, modQuery, "Compile", "No collection specified in Query")
	}
	if ins.Has("select") {
		//_ = "breakpoint"
		parm.Set("select", ins.Get("select", []string{}))
		//ret["select"] = ins.Get("select", M{}).(M)["select"]
		//ret["select"] = ins.Get("select", []string{})
	}
	if ins.Has("where") {
		parm.Set("find", ins.Get("where", M{}))
	}
	//_ = "breakpoint"
	command.Connection = q.Connection
	command.Settings = parm
	command.Text = tableName
	command.Type = commandType
	return command, nil
}

func (q *Query) Parse(qe *QE, ins M) interface{} {
	var v *QE
	result := M{}

	//-- field
	if qe.FieldOp == OpSelect {
		//_ = "breakpoint"
		return qe.Value
	} else

	//-- from
	if qe.FieldOp == OpFromTable {
		return qe.Value
	}

	//--- where
	if qe.FieldOp == OpEq {
		result.Set(qe.FieldId, qe.Value)
	} else if qe.FieldOp == OpNe {
		result.Set(qe.FieldId, M{}.Set("$ne", qe.Value))
	} else if qe.FieldOp == OpOr {
		ms := make([]M, 0)
		for _, v = range qe.Value.([]*QE) {
			ms = append(ms, q.Parse(v, ins).(M))
		}
		result.Set("$or", ms)
	} else if qe.FieldOp == OpAnd {
		ms := make([]M, 0)
		for _, v = range qe.Value.([]*QE) {
			ms = append(ms, q.Parse(v, ins).(M))
		}
		result.Set("$and", ms)
	}

	return result
}

package mongodb

import (
	//"fmt"
	. "github.com/eaciit/database/query"
	. "github.com/eaciit/toolkit"
)

type Query struct {
	QueryBase
	currentParseMode string
}

func (q *Query) Parse(qe QE, ins *M) interface{} {
	var v QE
	result := M{}

	if qe.FieldOp == OpEq {
		result.Set(qe.FieldId, qe.Value)
	} else if qe.FieldOp == OpNe {
		result.Set(qe.FieldId, M{}.Set("$ne", qe.Value))
	} else if qe.FieldOp == OpOr {
		ms := make([]M, 0)
		for _, v = range qe.Value.([]QE) {
			ms = append(ms, q.Parse(v, ins).(M))
		}
		result.Set("$or", ms)
	} else if qe.FieldOp == OpAnd {
		ms := make([]M, 0)
		for _, v = range qe.Value.([]QE) {
			ms = append(ms, q.Parse(v, ins).(M))
		}
		result.Set("$and", ms)
	}

	return result
}

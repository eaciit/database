package mongodb

import (
	//"fmt"
	"github.com/eaciit/database/query"
	. "github.com/eaciit/toolkit"
)

type MgoQuery struct {
	query.Query
	currentParseMode string
}

func (q *MgoQuery) Parse(result *M, ins *M, idx int) int {

	part := M{}
	command := M{}

	for _, v := range q.Elements {
		if v.FieldOp == query.OpOpenBracket {
			if part != nil {
				//command = append(command, part)
			}
			part = M{}
		} else if v.FieldOp == query.OpCloseBracket {
			//command = append(command, part)
		} else if v.FieldOp == query.OpOr {
			//part = part + " or "
		} else if v.FieldOp == query.OpAnd {
			//part = part + " and "
		} else if v.FieldOp == query.OpEq {
			p := new(M).Set("$eq", new(M).Set(v.FieldId, v.Value))
			command = p
			//command = append(command, p)
		}
	}

	result.Set("Data", command)
	return idx + 1
}

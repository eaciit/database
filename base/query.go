package base

import (
//"reflect"
//"strings"
)

// Query object is used to build database query command for both SQL and NoSQL
//
// select => q.Select(fields).From(table).Where(Eq("_id",1)).OrderBy(Order("Field1"),OrderDesc("Field2"))
// group => q.Aggregate(Sum("v11"),Sum("v2")).GroupBy("d1")
// insert => q.From(table).Insert()
// update => q.From(table).Update("field1","field2","field3")
// save => q.From(table).Save()
// delete => q.From(table).Delete()
// delete => q.From(table).Where(Eq("Status",2)).Delete()
//
// then invoke run. Run will generate cursor, updated value or error
//
// cursor, interface{}, error = q.Run(parms)

const (
	OpEq           = "$eq"
	OpNe           = "$ne"
	OpCommand      = "$command"
	OpGt           = "$gt"
	OpGte          = "$gte"
	OpLte          = "$lte"
	OpLt           = "$lt"
	OpBetween      = "$between"
	OpIn           = "$in"
	OpContains     = "$contains"
	OpStartWith    = "$startwith"
	OpEndWith      = "$endwith"
	OpAggregate    = "$aggregate"
	OpOpenBracket  = "$("
	OpCloseBracket = "$)"
	OpAnd          = "$and"
	OpOr           = "$or"
	OpChain        = "$chain"

	OpSelect    = "$select"
	OpSetfield  = "$setfield"
	OpFromTable = "$from"
	OpLimit     = "$limit"
	OpSkip      = "$skip"
	OpOrderBy   = "$order"
	OpGroupBy   = "$groupby"

	AggrSum    = "$sum"
	AggrCount  = "$count"
	AggrAvg    = "$average"
	AggrMin    = "$min"
	AggrMax    = "$max"
	AggrFirst  = "$first"
	AggrLast   = "$last"
	AggrMedian = "$median"
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

func Eq(field string, value interface{}) *QE {
	return &QE{field, OpEq, value}
}

func Ne(field string, value interface{}) *QE {
	return &QE{field, OpNe, value}

}

func Gt(field string, value interface{}) *QE {
	return &QE{field, OpGt, value}

}

func Gte(field string, value interface{}) *QE {
	return &QE{field, OpGte, value}

}

func Lt(field string, value interface{}) *QE {
	return &QE{field, OpLt, value}

}

func Lte(field string, value interface{}) *QE {
	return &QE{field, OpLte, value}

}

func Contains(field string, value ...string) *QE {
	return &QE{field, OpContains, value}

}

func StartWith(field string, value string) *QE {
	return &QE{field, OpStartWith, value}

}

func EndWith(field string, value string) *QE {
	return &QE{field, OpEndWith, value}

}

func Between(field string, from interface{}, to interface{}) *QE {
	return &QE{field, OpBetween, QRange{from, to}}
}

func In(field string, invalues ...interface{}) *QE {
	return &QE{field, OpIn, invalues}
}

func And(qes ...*QE) *QE {
	return &QE{"", OpAnd, qes}

}

func Or(qes ...*QE) *QE {
	return &QE{"", OpOr, qes}

}

//--- aggregate function
func Sum(field string) *QE {
	return &QE{field, AggrSum, nil}
}

func Count(field string) *QE {
	return &QE{field, AggrCount, nil}
}

func Avg(field string) *QE {
	return &QE{field, AggrAvg, nil}
}

func Min(field string) *QE {
	return &QE{field, AggrMin, nil}
}

func Max(field string) *QE {
	return &QE{field, AggrMax, nil}
}

func First(field string) *QE {
	return &QE{field, AggrFirst, nil}
}

func Last(field string) *QE {
	return &QE{field, AggrLast, nil}
}

func Median(field string) *QE {
	return &QE{field, AggrMedian, nil}
}

//--- creation function
func NewQuery(q IQuery) IQuery {
	q.SetQ(q)
	return q
}

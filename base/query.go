package base

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

	OpSelect    = "$select"
	OpFromTable = "$from"
	OpLimit     = "$limit"
	OpSkip      = "$skip"
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

func NewQuery(q IQuery) IQuery {
	q.SetQ(q)
	return q
}

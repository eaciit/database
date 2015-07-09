package base

/*
const (
	OrderDescending = "Desc"
	OrderAscending  = ""
	AggrSum         = "sum"
	AggrMax         = "max"
	AggrMin         = "min"
	AggrCount       = "count"
	AggrAvg         = "avg"
)

type IBuilder interface {
	RunAll(M) error
	Run(M) (ICursor, error)
	Adapter() IAdapters
	Select(...string) IBuilder
	From(string) IBuilder
	GroupBy(...string) IBuilder
	Aggregate(...AggregateItem) IBuilder
	OrderBy(...OrderItem) IBuilder
}

/  ******************
To run:
m := make([]M,0)
conn.Build().From("table").Select("field1","field2","field3").OrderBy("field1").RunAll(&m)
//conn.Build().From("table").Where(Or(W{"field1",EQ,""})).GroupBy("field1","field2","field3").Aggregate(A{AggrSum,"value1"})
******************   /

type A struct {
	AggregateOperation string
	FieldId            string
}

type O struct {
	FieldId       string
	SortDirection string
}

type W struct {
	FieldId string
	FieldOp string
	Value   interface{}
}

type BuilderBase struct {
	conn       IConnection
	table      string
	fields     []string
	groups     []string
	aggregates []AggregateItem
	orders     []OrderItem
}

func (b *BuilderBase) Init() {
	b.fields = make([]string, 0)
	b.groups = make([]string, 0)
	b.aggregates = make([]AggregateItem, 0)
	b.orders = make([]OrderItem, 0)
}

func (b *BuilderBase) Select(fs ...string) IBuilder {
	for _, v := range fs {
		b.fields = append(b.fields, v)
	}
	return b
}

func (b *BuilderBase) From(f string) IBuilder {
	b.table = f
	return b
}

func (b *BuilderBase) GroupBy(gs ...string) IBuilder {
	for _, v := range gs {
		b.groups = append(b.groups, g)
	}
	return b
}
*/

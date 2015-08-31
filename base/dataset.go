package base

type DataSet struct {
	Count int
	Data  []interface{}
	Model func() interface{}
}

func NewDataSet() *DataSet {
	ds := new(DataSet)
	ds.Data = make([]interface{}, 0)
	return ds
}

func (d *DataSet) SetModel(f func() interface{}) *DataSet {
	d.Model = f
	return d
}

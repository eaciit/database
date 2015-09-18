package mongodb

import (
	//"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/errorlib"
	"gopkg.in/mgo.v2"
)

type CursorType string

const (
	CursorType_Query CursorType = "query"
	CursorType_Pipe  CursorType = "pipe"
)

type Cursor struct {
	base.CursorBase
	Type     CursorType
	mgoSess  *mgo.Session
	mgoColl  *mgo.Collection
	mgoPipe  *mgo.Pipe
	mgoQuery *mgo.Query
	mgoIter  *mgo.Iter
}

func (c *Cursor) validate() error {
	if c.mgoQuery == nil && c.mgoPipe == nil {
		return errorlib.Error(packageName, modCursor, "validate", "Invalid Query or Pipe Object")
	}
	return nil
}

func (c *Cursor) FetchAll(result interface{}, closeCursor bool) error {
	var e error
	e = c.validate()
	if e != nil {
		return e
	}
	if c.Type == CursorType_Pipe {
		e = c.mgoPipe.All(result)
	} else {
		e = c.mgoQuery.All(result)
	}

	if closeCursor {
		c.Close()
	}
	return e
}

func (c *Cursor) ResetFetch() error {
	var e error
	e = c.validate()
	if e != nil {
		return e
	}
	if c.mgoIter != nil {
		c.mgoIter.Close()
	}

	if c.Type == CursorType_Pipe {
		c.mgoIter = c.mgoPipe.Iter()
	} else {
		c.mgoIter = c.mgoQuery.Iter()
	}

	return nil
}

func (c *Cursor) FetchClose(result interface{}) (bool, error) {
	defer c.Close()
	b, e := c.Fetch(result)
	return b, e
}

func (c *Cursor) Fetch(result interface{}) (bool, error) {
	var e error
	e = c.validate()
	if e != nil {
		return false, e
	}
	if c.mgoIter == nil {
		e = c.ResetFetch()
		if e != nil {
			return false, e
		}
	}
	boolIter := c.mgoIter.Next(result)

	return boolIter, nil
}

func (c *Cursor) FetchN(qty int, ds *base.DataSet, closeCursor bool) error {
	var e error
	e = c.validate()
	if e != nil {
		return e
	}
	if c.mgoIter == nil {
		e = c.ResetFetch()
		if e != nil {
			return e
		}
	}
	if closeCursor {
		defer c.Close()
	}

	scan := true
	i := 0
	for scan {
		dataHolder := ds.Model()
		if boolIter := c.mgoIter.Next(dataHolder); boolIter {
			ds.Data = append(ds.Data, dataHolder)
			i = i + 1
			if i == qty {
				scan = false
			}
		} else {
			scan = false
		}
	}

	return nil
}

func (c *Cursor) Count() int {
	var e error
	var n int
	if c.mgoQuery == nil {
		return 0
	}
	//_ = "breakpoint"
	if c.Type == CursorType_Pipe {
		return 0
	} else {
		n, e = c.mgoQuery.Count()
	}
	if e != nil {
		return 0
	}
	return n
}

func (c *Cursor) Close() {
	if c.mgoIter != nil {
		c.mgoIter.Close()
	}

	if c.mgoSess != nil && !c.Pooling() {
		c.mgoSess.Close()
	}
}

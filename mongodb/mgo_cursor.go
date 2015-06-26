package mongodb

import (
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

func (c *Cursor) FetchAll(result interface{}) error {
	var e error
	e = c.validate()
	if e != nil {
		return e
	}
	if c.Type == CursorType_Pipe {
		return c.mgoPipe.All(result)
	} else {
		return c.mgoQuery.All(result)
	}
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

func (c *Cursor) Count() int {
	if c.mgoQuery == nil {
		return 0
	}
	if c.Type == CursorType_Pipe {
		return 0
	} else {
		n, e := c.mgoQuery.Count()
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

	if c.mgoSess != nil {
		c.mgoSess.Close()
	}
}

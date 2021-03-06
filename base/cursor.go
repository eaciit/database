package base

import (
	"github.com/eaciit/errorlib"
)

type ICursor interface {
	ResetFetch() error
	Fetch(interface{}) (bool, error)
	FetchClose(interface{}) (bool, error)
	FetchN(int, *DataSet, bool) error
	FetchAll(interface{}, bool) error
	Count() int
	Close()
	GetQueryString() string

	Error() string
	SetError(string)
	SetPooling(bool)
	Pooling() bool
}

type CursorSourceType int

const (
	CursorTable CursorSourceType = iota
	CursorQuery
)

type CursorBase struct {
	CursorSource CursorSourceType
	Connection   IConnection
	QueryString  string

	errorTxt string
	pooling  bool
}

func (c *CursorBase) Pooling() bool {
	return c.pooling
}

func (c *CursorBase) SetPooling(p bool) {
	c.pooling = p
}

func (c *CursorBase) Error() string {
	return c.errorTxt
}

func (c *CursorBase) SetError(t string) {
	c.errorTxt = t
}

func (i *CursorBase) Fetch(result interface{}) (bool, error) {
	return false, errorlib.Error("database", "CursorBase", "Fetch", "Not yet implemented")
}

func (i *CursorBase) FetchClose(result interface{}) (bool, error) {
	return false, errorlib.Error("database", "CursorBase", "FetchClose", "Not yet implemented")
}

func (i *CursorBase) FetchN(nCount int, ds *DataSet, closeCursor bool) error {
	return errorlib.Error("database", "CursorBase", "FetchN", "Not yet implemented")
}

func (i *CursorBase) FetchAll(result interface{}, closeCursor bool) error {
	return errorlib.Error("database", "CursorBase", "FetchAll", "Not yet implemented")
}

func (i *CursorBase) ResetFetch() error {
	return errorlib.Error("database", "CursorBase", "ResetFetch", "Not yet implemented")
}

func (c *CursorBase) Count() int {
	return 0
}

func (c *CursorBase) Close() {
}

func (c *CursorBase) GetQueryString() string {
	return c.QueryString
}

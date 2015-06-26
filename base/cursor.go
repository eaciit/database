package base

import (
	"github.com/eaciit/errorlib"
)

type ICursor interface {
	ResetFetch() error
	Fetch(interface{}) (bool, error)
	FetchN(int, interface{}) (int, error)
	FetchAll(interface{}) error
	Count() int
	Close()
}

type CursorSourceType int

const (
	CursorTable CursorSourceType = iota
	CursorQuery
)

type CursorBase struct {
	CursorSource CursorSourceType
	Connection   IConnection
}

func (i *CursorBase) Fetch(result interface{}) (bool, error) {
	return false, errorlib.Error("database", "CursorBase", "Fetch", "Not yet implemented")
}

func (i *CursorBase) FetchN(nCount int, result interface{}) (int, error) {
	return 0, errorlib.Error("database", "CursorBase", "FetchN", "Not yet implemented")
}

func (i *CursorBase) FetchAll(result interface{}) error {
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

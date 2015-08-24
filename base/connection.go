package base

import (
	"errors"
	err "github.com/eaciit/errorlib"
	. "github.com/eaciit/toolkit"
)

type IConnection interface {
	Connect() error
	Execute(string, M) (int, error)
	Query() IQuery
	Table(string, M) ICursor
	//Adapter(string) IAdapter
	Close()
}

type ConnectionBase struct {
	Host     string
	UserName string
	Password string
	Database string
}

func (c *ConnectionBase) Execute(stmt string, parms M) (int, error) {
	return 0, err.Error(packageName, modConnection, "Execute", err.NotYetImplemented)
}

func (c *ConnectionBase) Connect() error {
	return errors.New("Method Connect is not yet implemented")
}

func (c *ConnectionBase) Close() {
}

/*
func (i *ConnectionBase) Adapter(tablename string) IAdapter {
	return new(AdapterBase)
}
*/

func (c *ConnectionBase) Query() IQuery {
	return NewQuery(new(QueryBase))
}

func (c *ConnectionBase) Table(tableName string, parms M) ICursor {
	return new(CursorBase)
}

package rdbms

import (
	"database/sql"
	"github.com/eaciit/database/base"
	err "github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
)

type Connection struct {
	*base.ConnectionBase
	ConnectionString string
	Sql              sql.DB
	Driver           string
}

func (c *Connection) Connect() error {
	sqlcon, e := sql.Open(c.Driver, c.ConnectionString)

	if e != nil {
		return err.Error(packageName, modConnection, "Connect", e.Error())
	}

	c.Sql = *sqlcon

	return nil
}

func (c *Connection) Query() base.IQuery {
	q := base.NewQuery(new(Query)).SetConnection(c).SetStringSign("")
	return q
}

func (c *Connection) Execute(stmt string, parms toolkit.M) (int, error) {
	return 0, nil
}

func (c *Connection) Command(cmdText string, settings map[string]interface{}) *base.CommandBase {
	return nil
}

func (c *Connection) Adapter(tableName string) base.IAdapter {
	return nil
}

func sel(q ...string) (r toolkit.M) {
	return nil
}

func (c *Connection) Table(tableName string, parms map[string]interface{}) base.ICursor {
	cs := new(Cursor)
	cs.QueryString = tableName
	cs.Connection = c
	return cs
}

func (c *Connection) Close() {
	c.Sql.Close()
}

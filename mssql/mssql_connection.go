package mssql

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/rdbms"
	"strings"
)

func NewConnection(host string, username string, password string, database string) base.IConnection {
	c := new(rdbms.Connection)
	c.Driver = "mssql"
	c.ConnectionBase = new(base.ConnectionBase)
	c.ConnectionBase.Host = host
	c.ConnectionBase.UserName = username
	c.ConnectionBase.Password = password
	c.ConnectionBase.Database = database
	c.ConnectionString = (func() string {
		if strings.Contains(host, ":") {
			host = fmt.Sprintf("%s;port=%s", strings.Split(host, ":")[0], strings.Split(host, ":")[1])
		}
		connectionString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", host, database, username, password)
		return connectionString
	}())

	return c
}

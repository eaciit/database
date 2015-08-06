package mysql

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/rdbms"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	ConnectionDefaultPort int = 3306
)

func NewConnection(host string, username string, password string, database string) base.IConnection {
	c := new(rdbms.Connection)
	c.Driver = "mysql"
	c.ConnectionBase = new(base.ConnectionBase)
	c.ConnectionBase.Host = host
	c.ConnectionBase.UserName = username
	c.ConnectionBase.Password = password
	c.ConnectionBase.Database = database
	c.ConnectionString = (func() string {
		if !strings.Contains(host, ":") {
			host = fmt.Sprintf("%s:%d", host, ConnectionDefaultPort)
		}
		connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, database)
		return connectionString
	}())

	return c
}

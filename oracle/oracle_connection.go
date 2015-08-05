package oracle

import (
	"database/sql"
	"fmt"
	"github.com/eaciit/database/base"
	err "github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	_ "github.com/mattn/go-oci8"
	"strconv"
	"strings"
)

const (
	ConnectionDefaultProtocol string = "TCP"
	ConnectionDefaultPort     int    = 1521
	ConnectionDefaultServer   string = "DEDICATED"
	ConnectionDefaultService  string = "ORCL"
)

type Connection struct {
	*base.ConnectionBase
	ConnectionString string
	Sql              sql.DB
}

func NewConnection(host string, username string, password string, database string) base.IConnection {
	c := new(Connection)
	c.ConnectionBase = new(base.ConnectionBase)
	c.ConnectionBase.Host = host
	c.ConnectionBase.UserName = username
	c.ConnectionBase.Password = password
	c.ConnectionBase.Database = database
	c.ConnectionString = (func() string {
		host, port := (func() (string, int) {
			port := ConnectionDefaultPort

			if strings.Contains(host, ":") {
				hostPort := strings.Split(host, ":")
				host = hostPort[0]
				port = (func() int {
					i, _ := strconv.Atoi(hostPort[1])
					return i
				}())
			}

			return host, port
		}())

		database, service := (func() (string, string) {
			service := ConnectionDefaultService

			if strings.Contains(database, "/") {
				databaseService := strings.Split(database, "/")
				database, service = databaseService[0], databaseService[1]
			}

			return database, service
		}())

		dsn := fmt.Sprintf(`%s/%s@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=%s)(HOST=%s)(PORT=%d)))(CONNECT_DATA=(SID=%s)(SERVER=%s)(SERVICE_NAME=%s)))`,
			username,
			password,
			ConnectionDefaultProtocol,
			host,
			port,
			database,
			ConnectionDefaultServer,
			service,
		)

		return dsn
	}())

	return c
}

func (c *Connection) Connect() error {
	sqlcon, e := sql.Open("oci8", c.ConnectionString)

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

func (c *Connection) Command(cmdText string, settings map[string]interface{}) *Command {
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

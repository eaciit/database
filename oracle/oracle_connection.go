package oracle

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/rdbms"
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

func NewConnection(host string, username string, password string, database string) base.IConnection {
	c := new(rdbms.Connection)
	c.Driver = "oci8"
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

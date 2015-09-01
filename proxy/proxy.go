package proxy

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/mongodb"
	"github.com/eaciit/database/mysql"
	"github.com/eaciit/database/oracle"
	err "github.com/eaciit/errorlib"
	"strings"
)

const (
	packageName = "database.proxy"
)

func NewConnection(connectionType string, host string, username string, password string, dbname string) (base.IConnection, error) {
	connectionType = strings.ToLower(connectionType)
	if connectionType == "mongodb" {
		c := mongodb.NewConnection(host, username, password, dbname)
		return c, nil
	} else if connectionType == "oracle" {
		c := oracle.NewConnection(host, username, password, dbname)
		return c, nil
	} else if connectionType == "mysql" {
		c := mysql.NewConnection(host, username, password, dbname)
		return c, nil
	}

	e := err.Error(packageName, "", "NewConnection", fmt.Sprintf("Connection type %s is not yet supported", connectionType))
	return nil, e
}

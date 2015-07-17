package mongodb

import (
	"fmt"
	. "github.com/eaciit/database/base"
	. "github.com/eaciit/toolkit"
	"testing"
)

var conn IConnection

func connect() {
	conn = NewConnection("localhost:27123", "", "", "ectest")
	e := conn.Connect()
	if e != nil {
		fmt.Println("Unable to connect " + e.Error())
	}
}

func prepareQuery() IQuery {
	qry := conn.Query()
	return qry.SetStringSign("\"").
		Select("fullname", "email").
		From("ORMUsers").
		Where(Or(Eq("_id", "user10"), Eq("_id", "user20")))
}

func TestQ(t *testing.T) {
	connect()
	defer conn.Close()
	_, e := prepareQuery().Build(nil)
	if e != nil {
		t.Error("Unable to parse Q :" + e.Error())
	} else {
		//fmt.Printf("Parse result: %v \n", GetJsonString(c))
	}
}

func TestR(t *testing.T) {
	connect()
	defer conn.Close()

	ms := make([]M, 0)
	prepareQuery().Cursor(nil).FetchAll(&ms, true)
	fmt.Printf("Result: \n%v\n", ms)
}

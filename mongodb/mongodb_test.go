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

func prepareCommand() (ICommand, error) {
	q := Eq("_id", "user01")

	qry := conn.Query()
	return qry.SetStringSign("\"").
		//Select("fullname", "email").
		From("ORMUsers").
		Where(q).
		Build(nil)
}

func TestQ(t *testing.T) {
	connect()
	defer conn.Close()
	_, e := prepareCommand()
	if e != nil {
		t.Error("Unable to parse Q :" + e.Error())
	} else {
		//fmt.Printf("Parse result: %v \n", GetJsonString(c))
	}
}

func TestR(t *testing.T) {
	connect()
	defer conn.Close()
	c, e := prepareCommand()

	if e != nil {
		t.Error("Unable to parse Q :" + e.Error())
	}

	ms := make([]M, 0)
	cursor, _, err := c.Run(nil, nil)
	if err != nil {
		t.Error("Unable to Execute command :" + e.Error())
	} else {
		fmt.Printf("Record found: %d \n", cursor.Count())
		cursor.FetchAll(&ms, true)
		fmt.Printf("Result: \n%v\n", ms)
	}
}

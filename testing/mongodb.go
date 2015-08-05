package main

import (
	"fmt"
	. "github.com/eaciit/database/base"
	. "github.com/eaciit/database/mongodb"
	. "github.com/eaciit/toolkit"
	_ "gopkg.in/mgo.v2/bson"
)

var conn IConnection

func connect() {
	if conn == nil {
		conn = NewConnection("192.168.0.200:27017", "", "", "ecshellx")
		e := conn.Connect()
		if e != nil {
			fmt.Println("Unable to connect " + e.Error())
		}
	}
}

func main() {
	connect()
	defer Close()

	ms := make([]M, 0)
	q := conn.Query().SetStringSign("\"").
		Select("RigType", "RigName").
		Limit(5).
		Where(Eq("WellName", "@val")).
		From("WEISWellActivities")
	c := q.Cursor(M{"@val": "WEISWellActivities"})
	e := c.FetchAll(&ms, true)
	if e != nil {
		fmt.Printf("Error fetch => %s \n", e.Error())
	}

	fmt.Println(ms)
}

func Close() {
	conn.Close()
	conn = nil
}

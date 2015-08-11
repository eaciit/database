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

	ms := []M{}
	q := conn.Query().
		//SetStringSign("\"").
		Select("RigType", "RigName", "WellName").
		From("WEISWellActivities").
		Where(Eq("WellName", "@val")).
		Limit(5)
	//c := q.Cursor(M{"@val": "Helix Q-4000"})
	c := q.Cursor(M{"@val": "PRINCESS P8"})
	fmt.Printf("Got %d record(s) \n", c.Count())
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

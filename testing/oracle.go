package main

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/oracle"
	"github.com/eaciit/toolkit"
)

var conn base.IConnection

func test(a interface{}) {
	*(a.(*string)) = "Asdfasdf"
}

func main() {
	conn := oracle.NewConnection("192.168.0.210:1521", "scott", "tiger", "ORCL/orcl.eaciit.local")
	conn.Connect()

	ms := []toolkit.M{}
	q := conn.Query().Select("customerid", "companyname").From("customers")
	c := q.Cursor(toolkit.M{"@uname": "noval", "@email": "noval@eaciit.com", "@action": "system"})
	c.FetchAll(&ms, true)

	fmt.Println("res", ms)
}

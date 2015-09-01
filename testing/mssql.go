package main

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/mssql"
	"github.com/eaciit/toolkit"
)

var conn base.IConnection

func main() {
	conn = mssql.NewConnection(`192.168.0.200`, "sa", "Password.Sql", "2015_anz_comtrade_source")
	conn.Connect()
	testSelectFromWhere()
	conn.Close()
}

func testSelectFromWhere() {
	q := conn.Query().
		SetStringSign("'").
		Select("country_id", "country_name").
		From("dim_country").
		Where(base.Lte("country_id", "@1"))
	c := q.Cursor(toolkit.M{"@1": 4})
	r := []toolkit.M{}
	e := c.FetchAll(&r, true)

	if e != nil {
		fmt.Println(e.Error())
	}

	fmt.Println("============== QUERY TEST SELECT-FROM-WHERE")
	fmt.Println(c.GetQueryString())

	for _, each := range r {
		fmt.Println(each)
	}
}

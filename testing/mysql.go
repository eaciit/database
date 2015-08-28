package main

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/mysql"
	"github.com/eaciit/toolkit"
)

var conn base.IConnection

func main() {
	conn = mysql.NewConnection("localhost", "root", "", "db_muslimorid")
	conn.Connect()
	// testInsertSet()
	testSelectFromWhereOrderLimitOffset()
	conn.Close()
}

func testInsert() {
	param := toolkit.M{"title": "tresno", "category": "cinta"}
	c, _, e := conn.Query().SetStringSign("'").Insert().From("tb_post").Run(param)

	if e != nil {
		fmt.Println(e.Error())
	}

	fmt.Println("============== QUERY TEST INSET-FROM-SET")
	fmt.Println(c.GetQueryString())
}

func testSelectFromWhereOrderLimitOffset() {
	q := conn.Query().
		SetStringSign("'").
		Select("id", "category", "author_name").
		From("tb_post").
		Where(base.Lte("id", "@1")).
		OrderBy("id asc").
		Limit(3).
		Skip(100)
	c := q.Cursor(toolkit.M{"@1": 373})
	r := []toolkit.M{}
	e := c.FetchAll(&r, true)

	if e != nil {
		fmt.Println(e.Error())
	}

	fmt.Println("============== QUERY TEST SELECT-FROM-WHERE-ORDERBY-LIMIT-OFFSET")
	fmt.Println(c.GetQueryString())

	for _, each := range r {
		fmt.Println(each)
	}
}

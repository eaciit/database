package main

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/oracle"
	"github.com/eaciit/toolkit"
)

var conn base.IConnection

func main() {
	conn = oracle.NewConnection("192.168.0.210:1521", "scott", "tiger", "ORCL/orcl.eaciit.local")
	conn.Connect()
	testSelectFromWhereOrder()
	testSelectFromLimitOffset()
	conn.Close()
}

func testSelectFromWhereOrder() {
	q := conn.Query().
		SetStringSign("'").
		Select("customerid", "companyname").
		From("customers").
		Where(
		base.Or(
			base.Eq("customerid", "@1"),
			base.Eq("customerid", "@2"),
			base.And(
				base.Eq("customerid", "@3"),
				base.Eq("companyname", "@4"))),
		base.Contains("companyname", "@5"),
		base.StartWith("companyname", "@6"),
		base.EndWith("companyname", "@7")).
		OrderBy("companyname asc", "customerid desc")
	c := q.Cursor(toolkit.M{
		"@1": "ANATR",
		"@2": "ANTON",
		"@3": "ALFKI",
		"@4": "Alfreds Futterkiste",
		"@5": "freds",
		"@6": "Alfreds",
		"@7": "Futterkiste",
	})
	r := []toolkit.M{}
	e := c.FetchAll(&r, true)

	if e != nil {
		fmt.Println(e.Error())
	}

	fmt.Println("============== QUERY TEST SELECT-FROM-WHERE-ORDERBY")
	fmt.Println(c.GetQueryString())

	for _, each := range r {
		fmt.Println(each)
	}
}

func testSelectFromLimitOffset() {
	q := conn.Query().
		SetStringSign("'").
		Select("customerid", "companyname").
		From("customers").
		Limit(4).
		Skip(2).
		OrderBy("companyname asc", "customerid desc")
	c := q.Cursor(toolkit.M{})
	r := []toolkit.M{}
	e := c.FetchAll(&r, true)

	if e != nil {
		fmt.Println(e.Error())
	}

	fmt.Println("============== QUERY TEST SELECT-FROM-LIMIT-OFFSET")
	fmt.Println(c.GetQueryString())

	for _, each := range r {
		fmt.Println(each)
	}
}

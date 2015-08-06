package main

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/mysql"
	"github.com/eaciit/toolkit"
)

var conn base.IConnection

func main() {
	conn := mysql.NewConnection("localhost", "root", "", "db_muslimorid")
	conn.Connect()

	ms := []toolkit.M{}
	q := conn.Query().
		SetStringSign("'").
		Select("id", "category", "author_name").
		From("tb_post").
		Where(
		base.Or(
			base.Eq("id", "@1"),
			base.Eq("id", "@2"))).
		OrderBy("id asc", "category desc")

	c := q.Cursor(toolkit.M{
		"@1": 375,
		"@2": 353,
	})

	e := c.FetchAll(&ms, true)
	if e != nil {
		fmt.Println(e.Error())
	}

	fmt.Println("res", ms)
}

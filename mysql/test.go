package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Row struct {
	id          int
	title       string
	category    string
	created_at  []uint8
	added_at    []uint8
	author_name string
	author_code string
	content     string
	link        string
	next_link   string
}

func main() {
	conn, err := sql.Open("mysql", "root@/db_muslimorid")

	if err != nil {
		fmt.Println(err.Error())
	}

	rows, err := conn.Query("select id, title, category, created_at, added_at, author_name, author_code, content, link, next_link from tb_post where id = ?", 375)

	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		row := Row{}
		column, _ := rows.Columns()

		c := []interface{}{&row.id, &row.title, &row.category, &row.created_at, &row.added_at, &row.author_name, &row.author_code, &row.content, &row.link, &row.next_link}

		fmt.Println(column)
		err = rows.Scan(c...)

		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(row.id, row.title)
	}
}

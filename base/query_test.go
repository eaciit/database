package base

import (
	"fmt"
	_ "github.com/eaciit/toolkit"
	"testing"
)

func TestQ(t *testing.T) {
	q := And(Or(
		And(Eq("username", "ariefdarmawan"), Eq("action", "system")),
		And(Eq("username", "someuser"), Eq("action", "log"))),
		Eq("username", "administrator"))

	qry := NewQuery(new(QueryBase))
	c, e := qry.SetStringSign("\"").Select("username", "action").From("UserLogs").Where(q).Build(nil)
	if e != nil {
		t.Error("Unable to parse Q :" + e.Error())
	} else {
		fmt.Printf("Parse result: %v \n", c)
	}
}

/*
qe := filter
q = new(new(QueryBase))
q.Select(selected).From(TableName).Where(qe).Command()
*/

/*
func TestQ_old(t *testing.T) {
	q := New(new(Query)).
		O().
		O().
		Eq("username", "ariefdarmawan").
		And().
		Eq("action", "system").
		C().
		Or().
		O().
		Eq("username", "someuser").
		And().
		Eq("action", "log").
		C().
		C().
		And().Chain(New(new(Query)).Eq("username", "administrator"))

	c := M{}
	e := q.Command(&c, nil)
	if e != nil {
		t.Error("Unable to parse Q")
	} else {
		fmt.Printf("Parse result: %v \n", c)
	}
}
*/

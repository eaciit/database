package query

import (
	"fmt"
	. "github.com/eaciit/toolkit"
	"testing"
)

func TestQ(t *testing.T) {
	q := And(Or(
		And(Eq("username", "ariefdarmawan"), Eq("action", "system")),
		And(Eq("username", "someuser"), Eq("action", "log"))),
		Eq("username", "administrator"))

	c := new(Result)
	e := new(Query).SetStringSign("\"").Command(c, nil, q)
	if e != nil {
		t.Error("Unable to parse Q")
	} else {
		fmt.Printf("Parse result: %v \n", c.Data)
	}
}

func Do(vs ...string) {
	fmt.Printf("You enter %d parms \n", len(vs))
	for _, v := range vs {
		fmt.Println(v)
	}
}

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

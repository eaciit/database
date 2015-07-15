package mongodb

import (
	"fmt"
	. "github.com/eaciit/database/query"
	. "github.com/eaciit/toolkit"
	"testing"
)

func TestQ(t *testing.T) {
	qes := And(Or(
		And(Eq("username", "ariefdarmawan"), Eq("action", "system")),
		And(Eq("username", "someuser"), Eq("action", "log"))),
		Eq("username", "administrator"))

	c := new(Result)
	q := new(MgoQuery)
	e := q.SetQ(q).SetStringSign("\"").Command(c, nil, qes)
	if e != nil {
		t.Error("Unable to parse Q")
	} else {
		fmt.Printf("Parse result: %v \n", GetJsonString(c.Data))
	}
}

func Do(vs ...string) {
	fmt.Printf("You enter %d parms \n", len(vs))
	for _, v := range vs {
		fmt.Println(v)
	}
}

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
	q := New(new(Query))
	e := q.SetQ(q).SetStringSign("\"").Command(c, nil, qes)
	if e != nil {
		t.Error("Unable to parse Q")
	} else {
		fmt.Printf("Parse result: %v \n", GetJsonString(c.Data))
	}
}

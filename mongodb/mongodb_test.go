package mongodb

import (
	"fmt"
	. "github.com/eaciit/database/query"
	. "github.com/eaciit/toolkit"
	"testing"
)

func TestQ(t *testing.T) {
	q := And(Or(
		And(Eq("username", "ariefdarmawan"), Eq("action", "system")),
		And(Eq("username", "someuser"), Eq("action", "log"))),
		Eq("username", "administrator"))

	c := new(Result)
	qry := New(new(Query))
	qry.SetStringSign("\"").Select("username", "action").From("UserLogs").Where(q).Build(c, nil)
	if c.Status != Status_OK {
		t.Error("Unable to parse Q :" + c.Message)
	} else {
		fmt.Printf("Parse result: %v \n", GetJsonString(c.Data))
	}
}

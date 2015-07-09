package query

import (
	"fmt"
	"testing"
	//. "github.com/eaciit/toolkit"
)

func TestQ(t *testing.T) {
	q := NewQuery(new(Query)).O().O().Eq("username", "ariefdarmawan").And().
		Eq("action", "system").C().Or().O().Eq("username", "someuser").And().Eq("action", "log").C().C().
		Or().Eq("username", "administrator")

	c := q.Parse(nil)
	if c == "" {
		t.Error("Unable to parse Q")
	} else {
		fmt.Printf("Parse:\n %v \n", c)
	}
}

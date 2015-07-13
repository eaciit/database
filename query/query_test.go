package query

import (
	"fmt"
	. "github.com/eaciit/toolkit"
	"testing"
)

func TestQ(t *testing.T) {
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
		fmt.Printf("Parse:\n %v \n", c)
	}
}

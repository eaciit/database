package query

import (
	"fmt"
	. "github.com/eaciit/toolkit"
	"testing"
<<<<<<< HEAD
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
=======
	"time"
)

func TestQ(t *testing.T) {
	q := new(Query).SetStringSign("'").Eq("field1", "Arief Darmawan").And().
		Eq("dateTrx", time.Now()).And().
		O().Eq("field2", 100).Or().Eq("field2", 200).C()

	cmd := q.Parse(nil)

	if cmd == "" {
>>>>>>> origin/master
		t.Error("Unable to parse Q")
	} else {
		fmt.Printf("Parse result: %v \n", cmd)
	//. "github.com/eaciit/toolkit"
)
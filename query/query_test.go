package query

import (
	"fmt"
	"testing"
	"time"
)

func TestQ(t *testing.T) {
	q := new(Query).SetStringSign("'").Eq("field1", "Arief Darmawan").And().
		Eq("dateTrx", time.Now()).And().
		O().Eq("field2", 100).Or().Eq("field2", 200).C()

	cmd := q.Parse(nil)

	if cmd == "" {
		t.Error("Unable to parse Q")
	} else {
		fmt.Printf("Parse result: %v \n", cmd)
	//. "github.com/eaciit/toolkit"
)
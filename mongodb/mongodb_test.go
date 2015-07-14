package mongodb

import (
	"fmt"
	"github.com/eaciit/database/query"
	. "github.com/eaciit/toolkit"
	"testing"
	"time"
)

func TestQ(t *testing.T) {
	q := query.New(new(MgoQuery)).SetStringSign("'").(*MgoQuery)
	q.Eq("field1", "Arief Darmawan")
	q.And()
	q.Eq("dateTrx", time.Now())
	q.And()
	q.O()
	q.Eq("field2", 100)
	q.Or()
	q.Eq("field2", 200)
	q.C()

	cmd := M{}
	e := q.Command(&cmd, &M{})

	if e != nil {
		t.Error("Unable to parse Q " + e.Error())
	} else {
		fmt.Printf("Parse result: %v \n", cmd)
	}
}

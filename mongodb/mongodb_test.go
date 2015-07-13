package mongodb

import (
	"fmt"
	. "github.com/eaciit/toolkit"
	"testing"
	"time"
)

func TestQ(t *testing.T) {
	Try(func() {
		q := new(MgoQuery).SetStringSign("'").(*MgoQuery)
		q.Eq("field1", "Arief Darmawan")
		q.And()
		q.Eq("dateTrx", time.Now())
		q.And()
		q.O()
		q.Eq("field2", 100)
		q.Or()
		q.Eq("field2", 200)
		q.C()

		cmd := q.Parse(nil)

		if cmd == "" {
			t.Error("Unable to parse Q")
		} else {
			fmt.Printf("Parse result: %v \n", cmd)
		}
	}).Catch(func(e interface{}) {
		//fmt.Println("Error at " + e.(error).Error())
		t.Error(e.(error).Error())
	}).Run()
}

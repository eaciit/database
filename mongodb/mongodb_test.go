package mongodb

import (
	"fmt"
	. "github.com/eaciit/database/base"
	. "github.com/eaciit/toolkit"
	_ "gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

var conn IConnection

func connect() {
	if conn == nil {
		conn = NewConnection("localhost:27123", "", "", "ecanadarko")
		_ = "breakpoint"
		e := conn.Connect()
		if e != nil {
			fmt.Println("Unable to connect " + e.Error())
		}
	}
}

func TestSimpleQuery(t *testing.T) {
	connect()
	defer Close()

	fmt.Print("Test Simple Query")
	ms := make([]M, 0)
	q := conn.Query().SetStringSign("\"").
		Select("wellinfos.wellid", "datereport.date").
		Limit(1).
		From("DataSummary")
	c := q.Cursor(M{"limit": 2, "skip": 10})
	e := c.FetchAll(&ms, true)
	if e != nil {
		t.Errorf("Error fetch => %s \n", e.Error())
	}

	if len(ms) == 0 {
		t.Errorf("Error on data => %v \n", ms)
	}
	fmt.Printf("... OK. Found %d records. Sample of first record is per below\n%v\n\n", len(ms), ms[0])
}

func TestAggregate(t *testing.T) {
	var e error
	connect()
	defer Close()

	fmt.Print("Test Aggregate")
	ms := []struct {
		Id struct {
			Datereport_Date  time.Time
			Wellinfos_WellId string
		} `bson:"_id"`
		Sum_TotalDays float64
		Sum_TotalCost float64
	}{}
	//ms := []M{}
	q := conn.Query().SetStringSign("\"").
		From("DataSummary").
		GroupBy("wellinfos.wellid", "datereport.date").
		Aggregate(Sum("totaldays"), Sum("totalcost"))
	e = q.Cursor(nil).FetchAll(&ms, true)
	if e != nil {
		t.Errorf("Error fetch => %s \n", e.Error())
	}
	fmt.Printf("... OK. Found %d records. Sample of first record is per below\n%v\n\n", len(ms), ms[0])
}

func TestQueryWithParam(t *testing.T) {
	connect()
	defer Close()

	fmt.Print("Test Query With Param")
	ms := make([]M, 0)
	q := conn.Query().
		//Select("_id", "wellinfos.wellid", "datereport.date").
		Where(Eq("wellinfos.wellid", "@wellid")).
		From("DataSummary")
	c := q.Cursor(M{"@wellid": "C5"})
	e := c.FetchAll(&ms, true)
	if e != nil {
		t.Errorf("Error fetch => %s \n", e.Error())
		return
	}

	if len(ms) == 0 {
		t.Errorf("Error on data => %v \n", ms)
		return
	}
	fmt.Printf("... OK. Found %d records. Sample of first record is per below\n%v\n\n", len(ms), ms[0])
}

func TestInsert(t *testing.T) {
	connect()
	defer Close()

	dt, _ := time.Parse("02-Jan-2006", "01-Apr-1980")
	data := M{"_id": "data01", "Title": "Ini adala data 1", "no1": 30,
		"tanggal": dt.UTC()}
	fmt.Print("Test Insert")
	q := conn.Query().From("TestTable").Insert()
	_, _, e := q.Run(M{"data": data})
	if e != nil {
		t.Errorf("Error: %s \n", e.Error())
		return
	} else {
		fmt.Printf("... OK\n")
	}
}

func TestUpdate(t *testing.T) {
	connect()
	defer Close()

	data := M{"_id": "data01", "Title": "Ini adala data 1 - Yg di Update", "no1": 30, "tanggal": MakeDate("02-Jan-2006", "01-Apr-1980").UTC()}
	fmt.Print("Test Update")
	q := conn.Query().From("TestTable").Update().SetFields("tanggal")
	_, _, e := q.Run(M{"data": data})
	if e != nil {
		t.Errorf("Error: %s \n", e.Error())
		return
	} else {
		fmt.Printf("... OK\n")
	}
}

func TestDelete(t *testing.T) {
	connect()
	defer Close()

	fmt.Print("Test Delete	")
	datas := []M{}
	qget := conn.Query().From("TestTable").Where(Eq("_id", "@id"))
	e := qget.Cursor(M{"@id": "data01"}).FetchAll(&datas, true)
	if e != nil {
		t.Errorf("Unable to load data = %s \n", e.Error())
		return
	}

	if len(datas) == 0 {
		t.Errorf("No data found match the criteria \n")
		return
	} else {
		fmt.Print("Test Delete")
		q := conn.Query().From("TestTable").Delete()
		_, _, e = q.Run(M{"data": datas[0]})
		if e != nil {
			t.Errorf("Error: %s \n", e.Error())
			return
		} else {
			fmt.Printf("... OK\n")
		}
	}
}

func Close() {
	conn.Close()
	conn = nil
}

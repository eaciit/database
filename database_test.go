package database

import (
	"fmt"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/mongodb"
	"strings"
	"testing"
)

func connect() (base.IConnection, error) {
	conn := mongodb.NewConnection("localhost:27017", "", "", "ectest")
	e := conn.Connect()
	return conn, e
}

type UserModel struct {
	Id       string `bson:"_id"`
	FullName string
	Enable   int
}

func TestDbConnect(t *testing.T) {
	fmt.Println("Test Connection to Database")
	fmt.Println("===========================")
	conn, e := connect()
	if e != nil {
		t.Error(fmt.Sprintf("%s \n", e.Error()))
		return
	}
	fmt.Println("Connected !")
	fmt.Println("")
	defer conn.Close()
}

/*
func TestDbInsert(t *testing.T) {
	fmt.Println("Test Insert Data to Database")
	fmt.Println("===========================")
	conn, e := connect()
	if e != nil {
		t.Error(fmt.Sprintf("%s \n", e.Error()))
		return
	}
	user := UserModel{"user02", "User 02A", 1}
	_, e = conn.Execute("Users", bson.M{"operation": "delete", "data": &user, "find": bson.M{"_id": "user02"}})
	if e != nil {
		t.Error(fmt.Sprintf("%s \n", e.Error()))
		return
	}
	fmt.Println("Inserted !")
	fmt.Println("")
	defer conn.Close()
}
*/

func TestAdapter(t *testing.T) {
	conn, e := connect()
	if e != nil {
		t.Error(fmt.Sprintf("%s \n", e.Error()))
		return
	}
	defer conn.Close()

	a := conn.Adapter("Users")
	cursor, _, e := a.Run(base.DB_SELECT, nil, nil)
	if e != nil {
		t.Error(fmt.Sprintf("%s \n", e.Error()))
		return
	}
	defer cursor.Close()
	u := new(UserModel)
	ok, _ := cursor.Fetch(u)
	if ok {
		fmt.Println(u)

		cursor.ResetFetch()
		idx := 0
		for ok, _ = cursor.Fetch(u); ok == true; ok, _ = cursor.Fetch(u) {
			idx = idx + 1

			if strings.Contains(u.FullName, "1") == true {
				u.FullName = "Employee 1 Name"
				a.Run(base.DB_SAVE, u, nil)
			}

			if strings.Contains(u.FullName, "Name") == false {
				u.FullName = u.Id + "'s Name"
				a.Run(base.DB_UPDATE, u, nil)
			}
			fmt.Printf("%d => %v \n", idx, u)
		}
	}
	defer cursor.Close()
}

/*
func TestQuery(t *testing.T) {
	fmt.Println("Test Run Query")
	fmt.Println("===========================")

	conn, e := connect()
	if e != nil {
		t.Error(fmt.Sprintf("%s \n", e.Error()))
		return
	}
	defer conn.Close()

	cursor := conn.Table("Users", bson.M{})
	if e != nil {
		t.Error(fmt.Sprintf("%s \n", e.Error()))
		return
	}
	defer cursor.Close()
	fmt.Printf("Number of users found: %d\n", cursor.Count())

	doc := new(UserModel)
	b := true
	idx := 0
	for b == true {
		idx++
		b, e = cursor.Fetch(doc)
		if b {
			if e != nil {
				t.Error(fmt.Sprintf("%s \n", e.Error()))
				//return
			} else {
				fmt.Printf("User Info %d: ID:%s FN:%s\n", idx, doc.Id, doc.FullName)
			}
		}
	}

	var result []UserModel
	e = cursor.FetchAll(&result)
	if e != nil {
		t.Error(fmt.Sprintf("%s \n", e.Error()))
		return
	}
	fmt.Printf("User Info:\n%v\n", result)

	fmt.Println("OK !")
	fmt.Println("")
}
*/

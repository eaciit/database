package mongodb

import (
	"github.com/eaciit/database/base"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	//"fmt"
)

type Adapter struct {
	base.AdapterBase
	mgoColl *mgo.Collection
}

type Command struct {
	base.CommandBase
}

func (c *Command) Run(data interface{}, parms map[string]interface{}) (base.ICursor, int, error) {
	var e error
	sess, mgoColl := c.Connection.(*Connection).CopySession(c.Text)
	defer sess.Close()

	_ = "breakpoint"
	var find bson.M
	if c.Type != base.DB_SELECT {
		var idField reflect.Value
		idField, hasField := toolkit.GetField(data, "Id")
		if hasField {
			find = bson.M{"_id": idField.Interface()}
			/*
				if idField.Kind() == Int {
					find := bson.M{"_id": idField.Int()}
				} else if idField.Kind() == String {
					find := bson.M{"_id": idField.String()}
				} else {
					find := bson.M{"_id": idField.Interface()}
				}
			*/
		}
	} else {
		find, hasFind := parms["find"]
		sort, hasSort := parms["sort"]
		skip, hasSkip := parms["skip"]
		limit, hasLimit := parms["limit"]

		var cursorParm bson.M
		if hasFind {
			cursorParm = bson.M{"find": find}
		}

		if hasSort {
			cursorParm["sort"] = sort
		}

		if hasSkip {
			cursorParm["skip"] = skip
		}

		if hasLimit {
			cursorParm["limit"] = limit
		}

		cursor := c.Connection.Table(c.Text, cursorParm)
		return cursor, 0, nil
	}

	if c.Type == base.DB_INSERT {
		e = mgoColl.Insert(data)
	} else if c.Type == base.DB_UPDATE {
		e = mgoColl.Update(find, data)
	} else if c.Type == base.DB_DELETE {
		e = mgoColl.Remove(find)
	} else if c.Type == base.DB_SAVE {
		_ = "breakpoint"
		_, e = mgoColl.Upsert(find, data)
		if e == nil {
			return nil, 0, nil
		}
	}
	if e != nil {
		return nil, 0, errorlib.Error(packageName, modCommand+"."+c.Type, "Run", e.Error())
	}
	return nil, 0, nil
}

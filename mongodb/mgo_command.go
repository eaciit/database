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

func (c *Command) Run(data interface{}, parms toolkit.M) (base.ICursor, int, error) {
	var e error
	//_ = "breakpoint"
	sess, mgoColl := c.Connection.(*Connection).CopySession(c.Text)
	defer sess.Close()

	var find bson.M
	if c.Type != base.DB_SELECT {
		var idField reflect.Value
		idField, hasIdField := toolkit.Field(data, "Id")
		if hasIdField {
			_ = bson.M{"_id": idField.Interface()}
		} else {
			_, _ = parms["find"]
		}
	} else {
		find, hasFind := parms["find"]
		sort, hasSort := parms["sort"]
		skip, hasSkip := parms["skip"]
		fields, hasFields := parms["select"]
		limit, hasLimit := parms["limit"]

		cursorParm := bson.M{}
		if c.Settings.Has("find") {
			cursorParm["find"] = c.Settings.Get("find", nil)
		}
		if hasFind {
			if c.Settings.Has("find") {
				cursorParm["find"] = bson.M{"$and": []interface{}{c.Settings["find"], find}}
			} else {
				cursorParm["find"] = find
			}
		}

		if hasFields {
			cursorParm["select"] = fields
		} else {
			if c.Settings.Has("select") {
				cursorParm["select"] = c.Settings.Get("select", nil)
			}
		}

		if hasSort {
			cursorParm["sort"] = sort
		} else {
			if c.Settings.Has("sort") {
				cursorParm["sort"] = c.Settings.Get("sort", nil)
			}
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
		//_ = "breakpoint"
		_, e = mgoColl.Upsert(find, data)
		if e == nil {
			return nil, 0, nil
		}
	}
	if e != nil {
		return nil, 0, errorlib.Error(packageName, modCommand+"."+string(c.Type), "Run", e.Error())
	}
	return nil, 0, nil
}

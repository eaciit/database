package mongodb

import (
	//"fmt"
	db "github.com/eaciit/database/base"
	err "github.com/eaciit/errorlib"
	"gopkg.in/mgo.v2"
	. "gopkg.in/mgo.v2/bson"
)

type Connection struct {
	*db.ConnectionBase
	mses *mgo.Session
	mdb  *mgo.Database
}

func NewConnection(host string, username string, password string, database string) db.IConnection {
	c := new(Connection)
	c.ConnectionBase = new(db.ConnectionBase)
	c.ConnectionBase.Host = host
	c.ConnectionBase.UserName = username
	c.ConnectionBase.Password = password
	c.ConnectionBase.Database = database
	return c
}

func (c *Connection) CopySession(tableName string) (*mgo.Session, *mgo.Collection) {
	copySession := c.mses.Copy()
	coll := copySession.DB(c.Database).C(tableName)
	return copySession, coll
}

func (c *Connection) Connect() error {
	sess, e := mgo.Dial(c.Host)
	if e != nil {
		return err.Error(packageName, modConnection, "Connect", e.Error())
	}
	sess.SetMode(mgo.Monotonic, true)
	//mdb := sess.DB(c.Database)
	c.mses = sess
	//c.mdb = mdb
	return nil
}

func (c *Connection) Query() db.IQuery {
	//_ = "breakpoint"
	q := db.NewQuery(new(Query))
	q.SetConnection(c)
	return q
}

func (c *Connection) Execute(stmt string, parms map[string]interface{}) (int, error) {
	var e error
	sess, coll := c.CopySession(stmt)
	defer sess.Close()

	//coll := c.mdb.C(stmt)
	op := ""
	ok := true
	if val, ok := parms["operation"]; !ok {
		return 0, err.Error(packageName, modConnection, "Execute", "Invalid operation in parms")
	} else {
		op = val.(string)
	}

	var data interface{}
	if data, ok = parms["data"]; !ok {
		return 0, err.Error(packageName, modConnection, "Execute", "Data is not valid")
	}

	if op == "insert" {
		e = coll.Insert(data)
		if e != nil {
			return 0, err.Error(packageName, modConnection, "Execute - Insert", e.Error())
		}
	} else {
		find, _ := parms["find"]
		if op == "update" {
			coll.Update(find, data)
		} else if op == "delete" {
			coll.Remove(find)
		}
	}
	return 0, nil
}

func (c *Connection) Command(cmdText string, settings map[string]interface{}) *Command {
	cmd := new(Command)
	cmd.Connection = c
	cmd.Text = cmdText
	cmd.Settings = settings
	return cmd
}

func (c *Connection) Adapter(tableName string) db.IAdapter {
	a := new(Adapter)
	a.Connection = c
	//a.mgoColl = c.mdb.C(tableName)
	a.SetCommand(db.DB_INSERT, c.Command(tableName, nil))
	a.SetCommand(db.DB_UPDATE, c.Command(tableName, nil))
	a.SetCommand(db.DB_DELETE, c.Command(tableName, nil))
	a.SetCommand(db.DB_SELECT, c.Command(tableName, nil))
	a.SetCommand(db.DB_SAVE, c.Command(tableName, nil))
	return a
}

func sel(q ...string) (r M) {
	r = make(M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}

func (c *Connection) Table(tableName string, parms map[string]interface{}) db.ICursor {
	cs := new(Cursor)
	cs.CursorSource = db.CursorTable
	cs.Connection = c
	cs.mgoSess, cs.mgoColl = c.CopySession(tableName)

	pipe, hasPipe := parms["pipe"]
	find, hasFind := parms["find"]
	sort, hasSort := parms["sort"]
	skip, hasSkip := parms["skip"]
	selectFields, hasSelectFields := parms["select"]
	limit, hasLimit := parms["limit"]

	if hasPipe {
		cs.mgoPipe = cs.mgoColl.Pipe(pipe).AllowDiskUse()
		cs.Type = CursorType_Pipe
	} else {
		cs.Type = CursorType_Query
		if hasFind {
			cs.mgoQuery = cs.mgoColl.Find(find)
		} else {
			cs.mgoQuery = cs.mgoColl.Find(nil)
		}

		if hasSelectFields {
			selecteds := sel(selectFields.([]string)...)
			cs.mgoQuery = cs.mgoQuery.Select(selecteds)
		}

		if hasSort {
			cs.mgoQuery = cs.mgoQuery.Sort(sort.(string))
		}
		if hasSkip {
			cs.mgoQuery = cs.mgoQuery.Skip(skip.(int))
		}
		if hasLimit {
			cs.mgoQuery = cs.mgoQuery.Limit(limit.(int))
		}
	}
	return cs
}

func (c *Connection) Close() {
	if c.mses != nil {
		c.mses.Close()
	}
}

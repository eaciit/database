package mongodb

import (
	"encoding/json"
	"fmt"
	. "github.com/eaciit/database/base"
	"github.com/eaciit/errorlib"
	. "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2"
	"os"
	"strings"
)

type Query struct {
	QueryBase
	currentParseMode string
}

func (q *Query) Parse(qe *QE, ins M) interface{} {
	var v *QE
	result := M{}
	//fmt.Printf("Debug m: %s \n", ins)

	//-- field
	if qe.FieldOp == OpSelect {
		////_ = "breakpoint"
		return qe.Value
	} else if qe.FieldOp == OpSetfield {
		return qe.Value
	} else

	//-- from
	if qe.FieldOp == OpFromTable {
		return qe.Value
	}

	//--- where
	if qe.FieldOp == OpEq {
		//_ = "breakpoint"
		value := q.ParseValue(qe.Value, ins)
		result.Set(qe.FieldId, value)
	} else if qe.FieldOp == OpNe {
		result.Set(qe.FieldId, M{}.Set("$ne", q.ParseValue(qe.Value, ins)))
	} else if qe.FieldOp == OpGt {
		result.Set(qe.FieldId, M{}.Set("$gt", q.ParseValue(qe.Value, ins)))
	} else if qe.FieldOp == OpGte {
		result.Set(qe.FieldId, M{}.Set("$gte", q.ParseValue(qe.Value, ins)))
	} else if qe.FieldOp == OpLt {
		result.Set(qe.FieldId, M{}.Set("$lt", q.ParseValue(qe.Value, ins)))
	} else if qe.FieldOp == OpLte {
		result.Set(qe.FieldId, M{}.Set("$lte", q.ParseValue(qe.Value, ins)))
	} else if qe.FieldOp == OpIn {
		result.Set(qe.FieldId, M{}.Set("$in", qe.Value))
		//fmt.Printf("value:%v\nresult:\n%v\n", JsonString(qe.Value), JsonString(result))
	} else if qe.FieldOp == OpContains {
		result.Set(qe.FieldId, M{}.
			Set("$regex", fmt.Sprintf(".*%s.*", qe.Value.([]string)[0])).
			Set("$options", "i"))
		//fmt.Printf("value:%v\nresult:\n%v\n", JsonString(qe.Value), JsonString(result))
	} else if qe.FieldOp == OpNin {
		result.Set(qe.FieldId, M{}.Set("$nin", qe.Value))
	} else

	//--- sort
	if qe.FieldOp == OpOrderBy {
		return qe.Value
	} else

	//--- Aggregate
	if qe.FieldOp == OpFlatten {
		return qe.Value
	} else if qe.FieldOp == OpGroupBy {
		groups := M{}
		groupValues := qe.Value.([]string)
		for _, v := range groupValues {
			groups.Set(strings.Replace(v, ".", "_", -1), "$"+v)
		}
		return groups
	} else if qe.FieldOp == OpAggregate {
		aggrs := M{}
		aggrQes := qe.Value.([]*QE)
		for _, qe := range aggrQes {
			m_aggr := M{}
			if qe.FieldOp == AggrSum {
				m_aggr.Set("$sum", "$"+qe.FieldId)
			} else if qe.FieldOp == AggrCount {
				m_aggr.Set("$sum", 1)
			} else if qe.FieldOp == AggrAvg {
				m_aggr.Set("$avg", "$"+qe.FieldId)
			} else if qe.FieldOp == AggrMin {
				m_aggr.Set("$min", "$"+qe.FieldId)
			} else if qe.FieldOp == AggrMax {
				m_aggr.Set("$max", "$"+qe.FieldId)
			} else if qe.FieldOp == AggrFirst {
				m_aggr.Set("$first", "$"+qe.FieldId)
			} else if qe.FieldOp == AggrLast {
				m_aggr.Set("$last", "$"+qe.FieldId)
			}
			aggrs.Set(strings.Replace(qe.FieldOp, "$", "", -1)+"_"+
				strings.Replace(qe.FieldId, ".", "_", -1), m_aggr)
		}
		////_ = "breakpoint"
		return aggrs
	} else
	//--- Skip & Limit
	if qe.FieldOp == OpSkip {
		return qe.Value
	} else if qe.FieldOp == OpLimit {
		return qe.Value
	} else
	//--- And
	if qe.FieldOp == OpOr {
		ms := make([]M, 0)
		for _, v = range qe.Value.([]*QE) {
			ms = append(ms, q.Parse(v, ins).(M))
		}
		result.Set("$or", ms)
	} else if qe.FieldOp == OpAnd {
		ms := make([]M, 0)
		for _, v = range qe.Value.([]*QE) {
			ms = append(ms, q.Parse(v, ins).(M))
			//_ = "breakpoint"
		}
		result.Set("$and", ms)
	} else if qe.FieldOp == OpWhereString {
		return qe.Value.(string)
	} else {
		return nil
	}
	return result
}

func mapEither(m1 M, m2 M, element string) (interface{}, bool) {
	if m1.Has(element) && m1[element] != nil {
		return m1[element], true
	} else if m2.Has(element) && m2[element] != nil {
		return m2[element], true
	} else {
		return nil, false
	}
}

func (q *Query) ParseWhereString(whereInterface interface{}) interface{} {
	var res interface{}
	err := json.Unmarshal([]byte(whereInterface.(string)), &res)

	if err != nil {
		panic(err.Error())
		os.Exit(0)
	}

	return res
}

func (q *Query) Compile(ins M) (ICursor, interface{}, error) {
	var e error
	s := q.Settings()
	//fmt.Printf("Settings %v \n", s)
	tableName := s.Get("from", "").([]string)[0]
	if tableName == "" {
		return nil, nil, errorlib.Error(packageName, modQuery, "Run",
			fmt.Sprint("No table / data source name specified"))
	}

	if ins == nil {
		ins = M{}
	}

	commandType := q.CommandType(s)

	var find M
	data, hasData := ins["data"]
	findQE, hasFind := ins["find"]
	if hasFind {
		find = q.Parse(findQE.(*QE), nil).(M)
	}
	if commandType == DB_SELECT {
		// read setting from main parms
		sort, hasSort := mapEither(ins, s, "sort")
		if hasSort == false {
			sort, hasSort = mapEither(ins, s, "orderby")
		}
		skip, hasSkip := mapEither(ins, s, "skip")
		fields, hasFields := mapEither(ins, s, "select")
		limit, hasLimit := mapEither(ins, s, "limit")

		// if not available then read it from settings
		cursorParm := M{}
		////_ = "breakpoint"
		if s.Has("where") {
			//fmt.Println("...", s.Get("where", nil))
			cursorParm["find"] = s.Get("where", nil)
		}
		if s.Has("whereString") {
			cursorParm["find"] = q.ParseWhereString(s.Get("whereString", "{}"))
		}
		if hasFind {
			if s.Has("find") {
				cursorParm["find"] = M{"$and": []interface{}{s["where"], find}}
			} else {
				cursorParm["find"] = find
			}
		}

		if hasFields {
			cursorParm["select"] = fields
		} else {
			if s.Has("select") {
				cursorParm["select"] = s.Get("select", nil)
			}
		}

		if hasSort {
			cursorParm["sort"] = sort
		} else {
			if s.Has("sort") {
				cursorParm["sort"] = s.Get("sort", nil)
			}
		}

		if hasSkip {
			cursorParm["skip"] = skip
		}

		if hasLimit {
			cursorParm["limit"] = limit
		}

		//--- handle aggregat compilation
		aggrs, hasaggr := s["aggregate"]
		groupby, hasgroupby := s["groupby"]
		if hasaggr || hasgroupby {
			pipes := []M{}
			pipe := M{}

			//--- check find
			if find, cursorHasFind := cursorParm["find"]; cursorHasFind {
				pipe.Set("$match", find)
				pipes = append(pipes, pipe)
			}

			//--- check flattern first
			flatten, hasflatten := s["flatten"]
			if hasflatten {
				pipe = M{}
				pipe.Set("$unwind", "$"+flatten.(string))
				pipes = append(pipes, pipe)
			}
			groupm := M{}

			groupm.Set("_id", groupby)
			pipe = M{}
			if hasaggr {
				for k, v := range aggrs.(M) {
					groupm.Set(k, v)
				}
			}
			pipe.Set("$group", groupm)
			pipes = append(pipes, pipe)

			if hasSkip {
				pipes = append(pipes, M{}.Set("$skip", cursorParm["skip"].(int)))
			}
			if hasLimit {
				pipes = append(pipes, M{}.Set("$limit", cursorParm["limit"].(int)))
			}

			cursorParm = M{}
			cursorParm.Set("pipe", pipes)
			//fmt.Printf("%v \n", JsonString(cursorParm))
		}
		//_ = "breakpoint"
		//fmt.Printf("Debug cp: %v \n", JsonString(cursorParm))
		if q.Pooling() {
			cursorParm.Set("pooling", q.Pooling())
		}
		cursor := q.Connection.Table(tableName, cursorParm)
		return cursor, 0, nil
	} else {
		var sess *mgo.Session
		var mgoColl *mgo.Collection
		if q.Pooling() == true {
			conn := q.Connection.(*Connection)
			mgoColl = conn.mses.DB(conn.Database).C(tableName)
		} else {
			sess, mgoColl = q.Connection.(*Connection).CopySession(tableName)
			defer sess.Close()
		}

		////_ = "breakpoint"
		hasIdField := false
		if hasData {
			////_ = "breakpoint"
			idField := Id(data)
			hasIdField = idField != nil
			if hasIdField {
				hasFind = true
				find = M{"_id": idField}
			}
		}

		multi := false
		if !hasFind {
			find = M{}
			multi = !hasIdField && commandType == DB_DELETE
		}

		if commandType == DB_INSERT {
			e = mgoColl.Insert(data)
		} else if commandType == DB_UPDATE {
			if multi {
				_, e = mgoColl.UpdateAll(find, data)
			} else {
				e = mgoColl.Update(find, data)
			}
		} else if commandType == DB_DELETE {
			if multi {
				_, e = mgoColl.RemoveAll(find)
			} else {
				e = mgoColl.Remove(find)
			}
		} else if commandType == DB_SAVE {
			////_ = "breakpoint"
			_, e = mgoColl.Upsert(find, data)
			if e == nil {
				return nil, 0, nil
			}
		}
		if e != nil {
			return nil, 0, errorlib.Error(packageName, modQuery+"."+string(commandType), "Run", e.Error())
		}
	}
	return nil, 0, nil
}

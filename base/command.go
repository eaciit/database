package base

import (
	//"github.com/eaciit/errorlib"
	. "github.com/eaciit/toolkit"
	"strings"
)

type ICommand interface {
	Prop(string, interface{})
	Run(interface{}, M) (ICursor, int, error)
}

type CommandBase struct {
	Connection IConnection
	Type       DB_OP
	Text       string
	Settings   M
}

func (c *CommandBase) Run(result interface{}, parms M) (ICursor, int, error) {
	return nil, 0, nil
}

func (c *CommandBase) Prop(fieldname string, result interface{}) {
	fieldname = strings.ToLower(fieldname)
	if fieldname == "type" {
		c.Type = result.(DB_OP)
	} else if fieldname == "text" {
		c.Text = result.(string)
	} else if fieldname == "settings" {
		c.Settings = result.(map[string]interface{})
	}
}
